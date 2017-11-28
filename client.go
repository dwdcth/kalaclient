package kala

import (
	"fmt"
	"gopkg.in/resty.v1"
	"net/http"
	"strings"
)

// KalaClient is the base struct for this package.
type KalaClient struct {
	apiEndpoint string
	requester   *resty.Request
}

// New is used to create a new KalaClient based off of the apiEndpoint
// Example:
// 		c := New("http://127.0.0.1:8000")
func New(apiEndpoint string) *KalaClient {
	if strings.HasSuffix(apiEndpoint, "/") {
		apiEndpoint = apiEndpoint[:len(apiEndpoint)-1]
	}

	return &KalaClient{
		apiEndpoint: apiEndpoint + ApiUrlPrefix,
		requester:   resty.SetHostURL(apiEndpoint + ApiUrlPrefix).SetRedirectPolicy(FlexibleRedirectPolicy(10)).R(),
	}
}

// CreateJob is used for creating a new job within Kala. It uses a map of
// strings to strings.
// Example:
// 		c := New("http://127.0.0.1:8000")
// 		body := map[string]string{
//			"schedule": "R2/2015-06-04T19:25:16.828696-07:00/PT10S",
//			"name":		"test_job",
//			"command": 	"bash -c 'date'",
//		}
//		id, err := c.CreateJob(body)
func (kc *KalaClient) CreateJob(body map[string]string) (string, error) {
	id := &AddJobResponse{}
	resp, err := kc.requester.SetBody(body).SetResult(id).Post(JobPath)
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != http.StatusCreated {
		return "", JobCreationError
	}
	return id.Id, nil
}

// GetJob is used to retrieve a Job from Kala by its ID.
// Example:
// 		c := New("http://127.0.0.1:8000")
//		id := "93b65499-b211-49ce-57e0-19e735cc5abd"
//		job, err := c.GetJob(id)
func (kc *KalaClient) GetJob(id string) (*Job, error) {
	j := &JobResponse{}
	resp, err := kc.requester.SetResult(j).Get(JobPath + id)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, JobNotFound
	}
	return j.Job, nil
}

// GetAllJobs returns a map of string (ID's) to job.Job's which contains
// all Jobs currently within Kala.
// Example:
// 		c := New("http://127.0.0.1:8000")
//		jobs, err := c.GetAllJobs()
func (kc *KalaClient) GetAllJobs() (map[string]*Job, error) {
	jobs := &ListJobsResponse{}
	resp, err := kc.requester.SetResult(jobs).Get(JobPath)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, GenericError
	}
	return jobs.Jobs, nil
}

// DeleteJob is used to delete a Job from Kala by its ID.
// Example:
// 		c := New("http://127.0.0.1:8000")
//		id := "93b65499-b211-49ce-57e0-19e735cc5abd"
//		ok, err := c.DeleteJob(id)
func (kc *KalaClient) DeleteJob(id string) (bool, error) {
	// nil is completely safe to use, as it is simply ignored in the sling library.
	resp, err := kc.requester.Delete(JobPath + id)
	if err != nil {
		return false, err
	}
	if resp.StatusCode() != http.StatusNoContent {
		return false, fmt.Errorf("Delete failed with a status code of %d", resp.StatusCode)
	}
	return true, nil
}

// GetJobStats is used to retrieve stats about a Job from Kala by its ID.
// Example:
// 		c := New("http://127.0.0.1:8000")
//		id := "93b65499-b211-49ce-57e0-19e735cc5abd"
//		stats, err := c.GetJobStats(id)
func (kc *KalaClient) GetJobStats(id string) ([]*JobStat, error) {
	js := &ListJobStatsResponse{}
	resp, err := kc.requester.SetResult(js).Get(JobPath + "stats/" + id)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, GenericError
	}
	return js.JobStats, nil
}

// StartJob is used to manually start a Job by its ID.
// Example:
// 		c := New("http://127.0.0.1:8000")
//		id := "93b65499-b211-49ce-57e0-19e735cc5abd"
//		ok, err := c.StartJob(id)
func (kc *KalaClient) StartJob(id string) (bool, error) {
	resp, err := kc.requester.Post(JobPath + "start/" + id)
	if err != nil {
		return false, err
	}
	if resp.StatusCode() != http.StatusNoContent {
		return false, nil
	}
	return true, nil
}

// GetKalaStats retrieves system-level metrics about Kala
// Example:
// 		c := New("http://127.0.0.1:8000")
//		stats, err := c.GetKalaStats()
func (kc *KalaClient) GetKalaStats() (*KalaStats, error) {
	ks := &KalaStatsResponse{}
	resp, err := kc.requester.SetResult(ks).Get("stats")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, GenericError
	}
	return ks.Stats, nil
}
