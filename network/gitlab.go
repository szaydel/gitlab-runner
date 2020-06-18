package network

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"gitlab.com/gitlab-org/gitlab-runner/common"
	"gitlab.com/gitlab-org/gitlab-runner/helpers"
)

const clientError = -100

var apiRequestStatuses = prometheus.NewDesc(
	"gitlab_runner_api_request_statuses_total",
	"The total number of api requests, partitioned by runner, endpoint and status.",
	[]string{"runner", "endpoint", "status"},
	nil,
)

type APIEndpoint string

const (
	APIEndpointRequestJob APIEndpoint = "request_job"
	APIEndpointUpdateJob  APIEndpoint = "update_job"
	APIEndpointPatchTrace APIEndpoint = "patch_trace"
)

type apiRequestStatusPermutation struct {
	runnerID string
	endpoint APIEndpoint
	status   int
}

type APIRequestStatusesMap struct {
	internal map[apiRequestStatusPermutation]int
	lock     sync.RWMutex
}

func (arspm *APIRequestStatusesMap) Append(runnerID string, endpoint APIEndpoint, status int) {
	arspm.lock.Lock()
	defer arspm.lock.Unlock()

	permutation := apiRequestStatusPermutation{runnerID: runnerID, endpoint: endpoint, status: status}

	if _, ok := arspm.internal[permutation]; !ok {
		arspm.internal[permutation] = 0
	}

	arspm.internal[permutation]++
}

// Describe implements prometheus.Collector.
func (arspm *APIRequestStatusesMap) Describe(ch chan<- *prometheus.Desc) {
	ch <- apiRequestStatuses
}

// Collect implements prometheus.Collector.
func (arspm *APIRequestStatusesMap) Collect(ch chan<- prometheus.Metric) {
	arspm.lock.RLock()
	defer arspm.lock.RUnlock()

	for permutation, count := range arspm.internal {
		ch <- prometheus.MustNewConstMetric(
			apiRequestStatuses,
			prometheus.CounterValue,
			float64(count),
			permutation.runnerID,
			string(permutation.endpoint),
			strconv.Itoa(permutation.status),
		)
	}
}

func NewAPIRequestStatusesMap() *APIRequestStatusesMap {
	return &APIRequestStatusesMap{
		internal: make(map[apiRequestStatusPermutation]int),
	}
}

type GitLabClient struct {
	clients map[string]*client
	lock    sync.Mutex

	requestsStatusesMap *APIRequestStatusesMap
}

func (n *GitLabClient) getClient(credentials requestCredentials) (c *client, err error) {
	n.lock.Lock()
	defer n.lock.Unlock()

	if n.clients == nil {
		n.clients = make(map[string]*client)
	}
	key := fmt.Sprintf(
		"%s_%s_%s_%s",
		credentials.GetURL(),
		credentials.GetToken(),
		credentials.GetTLSCAFile(),
		credentials.GetTLSCertFile(),
	)
	c = n.clients[key]
	if c == nil {
		c, err = newClient(credentials)
		if err != nil {
			return
		}
		n.clients[key] = c
	}

	return
}

func (n *GitLabClient) getLastUpdate(credentials requestCredentials) (lu string) {
	cli, err := n.getClient(credentials)
	if err != nil {
		return ""
	}
	return cli.getLastUpdate()
}

func (n *GitLabClient) getRunnerVersion(config common.RunnerConfig) common.VersionInfo {
	info := common.VersionInfo{
		Name:         common.NAME,
		Version:      common.VERSION,
		Revision:     common.REVISION,
		Platform:     runtime.GOOS,
		Architecture: runtime.GOARCH,
		Executor:     config.Executor,
		Shell:        config.Shell,
	}

	if executorProvider := common.GetExecutorProvider(config.Executor); executorProvider != nil {
		_ = executorProvider.GetFeatures(&info.Features)

		if info.Shell == "" {
			info.Shell = executorProvider.GetDefaultShell()
		}
	}

	if shell := common.GetShell(info.Shell); shell != nil {
		shell.GetFeatures(&info.Features)
	}

	return info
}

func (n *GitLabClient) doRaw(
	credentials requestCredentials,
	method, uri string,
	request io.Reader,
	requestType string,
	headers http.Header,
) (res *http.Response, err error) {
	c, err := n.getClient(credentials)
	if err != nil {
		return nil, err
	}

	return c.do(uri, method, request, requestType, headers)
}

func (n *GitLabClient) doJSON(
	credentials requestCredentials,
	method, uri string,
	statusCode int,
	request interface{},
	response interface{},
) (int, string, *http.Response) {
	c, err := n.getClient(credentials)
	if err != nil {
		return clientError, err.Error(), nil
	}

	return c.doJSON(uri, method, statusCode, request, response)
}

func (n *GitLabClient) getResponseTLSData(
	credentials requestCredentials,
	response *http.Response,
) (ResponseTLSData, error) {
	c, err := n.getClient(credentials)
	if err != nil {
		return ResponseTLSData{}, fmt.Errorf("couldn't get client: %w", err)
	}

	return c.getResponseTLSData(response.TLS)
}

func (n *GitLabClient) RegisterRunner(
	runner common.RunnerCredentials,
	parameters common.RegisterRunnerParameters,
) *common.RegisterRunnerResponse {
	// TODO: pass executor
	request := common.RegisterRunnerRequest{
		RegisterRunnerParameters: parameters,
		Token:                    runner.Token,
		Info:                     n.getRunnerVersion(common.RunnerConfig{}),
	}

	var response common.RegisterRunnerResponse
	result, statusText, resp := n.doJSON(
		&runner,
		http.MethodPost,
		"runners",
		http.StatusCreated,
		&request,
		&response,
	)
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}

	switch result {
	case http.StatusCreated:
		runner.Log().Println("Registering runner...", "succeeded")
		return &response
	case http.StatusForbidden:
		runner.Log().Errorln("Registering runner...", "forbidden (check registration token)")
		return nil
	case clientError:
		runner.Log().WithField("status", statusText).Errorln("Registering runner...", "error")
		return nil
	default:
		runner.Log().WithField("status", statusText).Errorln("Registering runner...", "failed")
		return nil
	}
}

func (n *GitLabClient) VerifyRunner(runner common.RunnerCredentials) bool {
	request := common.VerifyRunnerRequest{
		Token: runner.Token,
	}

	result, statusText, resp := n.doJSON(&runner, http.MethodPost, "runners/verify", http.StatusOK, &request, nil)
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}

	switch result {
	case http.StatusOK:
		// this is expected due to fact that we ask for non-existing job
		runner.Log().Println("Verifying runner...", "is alive")
		return true
	case http.StatusForbidden:
		runner.Log().Errorln("Verifying runner...", "is removed")
		return false
	case clientError:
		runner.Log().WithField("status", statusText).Errorln("Verifying runner...", "error")
		return true
	default:
		runner.Log().WithField("status", statusText).Errorln("Verifying runner...", "failed")
		return true
	}
}

func (n *GitLabClient) UnregisterRunner(runner common.RunnerCredentials) bool {
	request := common.UnregisterRunnerRequest{
		Token: runner.Token,
	}

	result, statusText, resp := n.doJSON(&runner, http.MethodDelete, "runners", http.StatusNoContent, &request, nil)
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}

	const baseLogText = "Unregistering runner from GitLab"
	switch result {
	case http.StatusNoContent:
		runner.Log().Println(baseLogText, "succeeded")
		return true
	case http.StatusForbidden:
		runner.Log().Errorln(baseLogText, "forbidden")
		return false
	case clientError:
		runner.Log().WithField("status", statusText).Errorln(baseLogText, "error")
		return false
	default:
		runner.Log().WithField("status", statusText).Errorln(baseLogText, "failed")
		return false
	}
}

func addTLSData(response *common.JobResponse, tlsData ResponseTLSData) {
	if tlsData.CAChain != "" {
		response.TLSCAChain = tlsData.CAChain
	}

	if tlsData.CertFile != "" && tlsData.KeyFile != "" {
		data, err := ioutil.ReadFile(tlsData.CertFile)
		if err == nil {
			response.TLSAuthCert = string(data)
		}
		data, err = ioutil.ReadFile(tlsData.KeyFile)
		if err == nil {
			response.TLSAuthKey = string(data)
		}
	}
}

func (n *GitLabClient) RequestJob(
	config common.RunnerConfig,
	sessionInfo *common.SessionInfo,
) (*common.JobResponse, bool) {
	request := common.JobRequest{
		Info:       n.getRunnerVersion(config),
		Token:      config.Token,
		LastUpdate: n.getLastUpdate(&config.RunnerCredentials),
		Session:    sessionInfo,
	}

	var response common.JobResponse
	result, statusText, httpResponse := n.doJSON(
		&config.RunnerCredentials,
		http.MethodPost,
		"jobs/request",
		http.StatusCreated,
		&request,
		&response,
	)

	n.requestsStatusesMap.Append(config.RunnerCredentials.ShortDescription(), APIEndpointRequestJob, result)

	switch result {
	case http.StatusCreated:
		config.Log().WithFields(logrus.Fields{
			"job":      response.ID,
			"repo_url": response.RepoCleanURL(),
		}).Println("Checking for jobs...", "received")

		tlsData, err := n.getResponseTLSData(&config.RunnerCredentials, httpResponse)
		if err != nil {
			config.Log().
				WithError(err).Errorln("Error on fetching TLS Data from API response...", "error")
		}
		addTLSData(&response, tlsData)

		return &response, true
	case http.StatusForbidden:
		config.Log().Errorln("Checking for jobs...", "forbidden")
		return nil, false
	case http.StatusNoContent:
		config.Log().Debugln("Checking for jobs...", "nothing")
		return nil, true
	case clientError:
		config.Log().WithField("status", statusText).Errorln("Checking for jobs...", "error")
		return nil, false
	default:
		config.Log().WithField("status", statusText).Warningln("Checking for jobs...", "failed")
		return nil, true
	}
}

func (n *GitLabClient) UpdateJob(
	config common.RunnerConfig,
	jobCredentials *common.JobCredentials,
	jobInfo common.UpdateJobInfo,
) common.UpdateState {
	request := common.UpdateJobRequest{
		Info:          n.getRunnerVersion(config),
		Token:         jobCredentials.Token,
		State:         jobInfo.State,
		FailureReason: jobInfo.FailureReason,
	}

	result, statusText, response := n.doJSON(
		&config.RunnerCredentials,
		http.MethodPut,
		fmt.Sprintf("jobs/%d", jobInfo.ID),
		http.StatusOK,
		&request,
		nil,
	)
	n.requestsStatusesMap.Append(config.RunnerCredentials.ShortDescription(), APIEndpointUpdateJob, result)

	remoteJobStateResponse := NewRemoteJobStateResponse(response)
	log := config.Log().WithFields(logrus.Fields{
		"code":       result,
		"job":        jobInfo.ID,
		"job-status": remoteJobStateResponse.RemoteState,
	})

	switch {
	case remoteJobStateResponse.IsAborted():
		log.Warningln("Submitting job to coordinator...", "aborted")
		return common.UpdateAbort
	case result == http.StatusOK:
		log.Debugln("Submitting job to coordinator...", "ok")
		return common.UpdateSucceeded
	case result == http.StatusNotFound:
		log.Warningln("Submitting job to coordinator...", "aborted")
		return common.UpdateAbort
	case result == http.StatusForbidden:
		log.WithField("status", statusText).Errorln("Submitting job to coordinator...", "forbidden")
		return common.UpdateAbort
	case result == clientError:
		log.WithField("status", statusText).Errorln("Submitting job to coordinator...", "error")
		return common.UpdateAbort
	default:
		log.WithField("status", statusText).Warningln("Submitting job to coordinator...", "failed")
		return common.UpdateFailed
	}
}

func (n *GitLabClient) PatchTrace(
	config common.RunnerConfig,
	jobCredentials *common.JobCredentials,
	content []byte,
	startOffset int,
) common.PatchTraceResult {
	id := jobCredentials.ID

	baseLog := config.Log().WithField("job", id)
	if len(content) == 0 {
		baseLog.Debugln("Appending trace to coordinator...", "skipped due to empty patch")
		return common.NewPatchTraceResult(startOffset, common.UpdateSucceeded, 0)
	}

	endOffset := startOffset + len(content)
	contentRange := fmt.Sprintf("%d-%d", startOffset, endOffset-1)

	headers := make(http.Header)
	headers.Set("Content-Range", contentRange)
	headers.Set("JOB-TOKEN", jobCredentials.Token)

	uri := fmt.Sprintf("jobs/%d/trace", id)
	request := bytes.NewReader(content)

	response, err := n.doRaw(&config.RunnerCredentials, "PATCH", uri, request, "text/plain", headers)
	if err != nil {
		config.Log().Errorln("Appending trace to coordinator...", "error", err.Error())
		return common.NewPatchTraceResult(startOffset, common.UpdateFailed, 0)
	}

	n.requestsStatusesMap.Append(
		config.RunnerCredentials.ShortDescription(),
		APIEndpointPatchTrace,
		response.StatusCode,
	)

	defer func() {
		_, _ = io.Copy(ioutil.Discard, response.Body)
		_ = response.Body.Close()
	}()

	tracePatchResponse := NewTracePatchResponse(response, baseLog)
	log := baseLog.WithFields(logrus.Fields{
		"sent-log":        contentRange,
		"job-log":         tracePatchResponse.RemoteRange,
		"job-status":      tracePatchResponse.RemoteState,
		"code":            response.StatusCode,
		"status":          response.Status,
		"update-interval": tracePatchResponse.RemoteTraceUpdateInterval,
	})

	return n.createPatchTraceResult(startOffset, tracePatchResponse, response, endOffset, log)
}

func (n *GitLabClient) createPatchTraceResult(
	startOffset int,
	tracePatchResponse *TracePatchResponse,
	response *http.Response,
	endOffset int,
	log *logrus.Entry,
) common.PatchTraceResult {
	result := common.PatchTraceResult{
		SentOffset:        startOffset,
		NewUpdateInterval: tracePatchResponse.RemoteTraceUpdateInterval,
	}

	switch {
	case tracePatchResponse.IsAborted():
		log.Warningln("Appending trace to coordinator...", "aborted")
		result.State = common.UpdateAbort

		return result

	case response.StatusCode == http.StatusAccepted:
		log.Debugln("Appending trace to coordinator...", "ok")
		result.SentOffset = endOffset
		result.State = common.UpdateSucceeded

		return result

	case response.StatusCode == http.StatusNotFound:
		log.Warningln("Appending trace to coordinator...", "not-found")
		result.State = common.UpdateNotFound

		return result

	case response.StatusCode == http.StatusRequestedRangeNotSatisfiable:
		log.Warningln("Appending trace to coordinator...", "range mismatch")
		result.SentOffset = tracePatchResponse.NewOffset()
		result.State = common.UpdateRangeMismatch

		return result

	case response.StatusCode == clientError:
		log.Errorln("Appending trace to coordinator...", "error")
		result.State = common.UpdateAbort

		return result

	default:
		log.Warningln("Appending trace to coordinator...", "failed")
		result.State = common.UpdateFailed

		return result
	}
}

func (n *GitLabClient) createArtifactsForm(mpw *multipart.Writer, reader io.Reader, baseName string) error {
	wr, err := mpw.CreateFormFile("file", baseName)
	if err != nil {
		return err
	}

	_, err = io.Copy(wr, reader)
	if err != nil {
		return err
	}
	return nil
}

func uploadRawArtifactsQuery(options common.ArtifactsOptions) url.Values {
	q := url.Values{}

	if options.ExpireIn != "" {
		q.Set("expire_in", options.ExpireIn)
	}

	if options.Format != "" {
		q.Set("artifact_format", string(options.Format))
	}

	if options.Type != "" {
		q.Set("artifact_type", options.Type)
	}

	return q
}

func (n *GitLabClient) UploadRawArtifacts(
	config common.JobCredentials,
	reader io.Reader,
	options common.ArtifactsOptions,
) common.UploadState {
	pr, pw := io.Pipe()
	defer func() { _ = pr.Close() }()

	mpw := multipart.NewWriter(pw)

	go func() {
		defer func() {
			_ = mpw.Close()
			_ = pw.Close()
		}()
		err := n.createArtifactsForm(mpw, reader, options.BaseName)
		if err != nil {
			_ = pw.CloseWithError(err)
		}
	}()

	query := uploadRawArtifactsQuery(options)

	headers := make(http.Header)
	headers.Set("JOB-TOKEN", config.Token)
	res, err := n.doRaw(
		&config,
		http.MethodPost,
		fmt.Sprintf("jobs/%d/artifacts?%s", config.ID, query.Encode()),
		pr,
		mpw.FormDataContentType(),
		headers,
	)

	log := logrus.WithFields(logrus.Fields{
		"id":    config.ID,
		"token": helpers.ShortenToken(config.Token),
	})

	if res != nil {
		log = log.WithField("responseStatus", res.Status)
	}

	messagePrefix := "Uploading artifacts to coordinator..."
	if options.Type != "" {
		messagePrefix = fmt.Sprintf("Uploading artifacts as %q to coordinator...", options.Type)
	}

	if err != nil {
		log.WithError(err).Errorln(messagePrefix, "error")
		return common.UploadFailed
	}
	defer func() {
		_, _ = io.Copy(ioutil.Discard, res.Body)
		_ = res.Body.Close()
	}()

	return n.determineUploadState(res.StatusCode, log, messagePrefix)
}

func (n *GitLabClient) determineUploadState(
	statusCode int,
	log *logrus.Entry,
	messagePrefix string,
) common.UploadState {
	switch statusCode {
	case http.StatusCreated:
		log.Println(messagePrefix, "ok")
		return common.UploadSucceeded
	case http.StatusForbidden:
		log.WithField("status", statusCode).Errorln(messagePrefix, "forbidden")
		return common.UploadForbidden
	case http.StatusRequestEntityTooLarge:
		log.WithField("status", statusCode).Errorln(messagePrefix, "too large archive")
		return common.UploadTooLarge
	case http.StatusServiceUnavailable:
		log.WithField("status", statusCode).Errorln(messagePrefix, "service unavailable")
		return common.UploadServiceUnavailable
	default:
		log.WithField("status", statusCode).Warningln(messagePrefix, "failed")
		return common.UploadFailed
	}
}

func (n *GitLabClient) DownloadArtifacts(
	config common.JobCredentials,
	artifactsFile string,
	directDownload *bool,
) common.DownloadState {
	query := url.Values{}

	if directDownload != nil {
		query.Set("direct_download", strconv.FormatBool(*directDownload))
	}

	headers := make(http.Header)
	headers.Set("JOB-TOKEN", config.Token)
	uri := fmt.Sprintf("jobs/%d/artifacts?%s", config.ID, query.Encode())

	res, err := n.doRaw(&config, http.MethodGet, uri, nil, "", headers)

	log := logrus.WithFields(logrus.Fields{
		"id":    config.ID,
		"token": helpers.ShortenToken(config.Token),
	})

	if res != nil {
		log = log.WithField("responseStatus", res.Status)
	}

	if err != nil {
		log.Errorln("Downloading artifacts from coordinator...", "error", err.Error())
		return common.DownloadFailed
	}
	defer func() {
		_, _ = io.Copy(ioutil.Discard, res.Body)
		_ = res.Body.Close()
	}()

	switch res.StatusCode {
	case http.StatusOK:
		file, err := os.Create(artifactsFile)
		if err == nil {
			defer func() { _ = file.Close() }()
			_, err = io.Copy(file, res.Body)
		}
		if err != nil {
			_ = file.Close()
			_ = os.Remove(file.Name())
			log.WithError(err).Errorln("Downloading artifacts from coordinator...", "error")
			return common.DownloadFailed
		}
		log.Println("Downloading artifacts from coordinator...", "ok")
		return common.DownloadSucceeded
	case http.StatusForbidden:
		log.WithField("status", res.Status).Errorln("Downloading artifacts from coordinator...", "forbidden")
		return common.DownloadForbidden
	case http.StatusNotFound:
		log.Errorln("Downloading artifacts from coordinator...", "not found")
		return common.DownloadNotFound
	default:
		log.WithField("status", res.Status).Warningln("Downloading artifacts from coordinator...", "failed")
		return common.DownloadFailed
	}
}

func (n *GitLabClient) ProcessJob(
	config common.RunnerConfig,
	jobCredentials *common.JobCredentials,
) (common.JobTrace, error) {
	trace, err := newJobTrace(n, config, jobCredentials)
	if err != nil {
		return nil, err
	}

	trace.start()
	return trace, nil
}

func NewGitLabClientWithRequestStatusesMap(rsMap *APIRequestStatusesMap) *GitLabClient {
	return &GitLabClient{
		requestsStatusesMap: rsMap,
	}
}

func NewGitLabClient() *GitLabClient {
	return NewGitLabClientWithRequestStatusesMap(NewAPIRequestStatusesMap())
}
