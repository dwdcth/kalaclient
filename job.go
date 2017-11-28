package kala

import (
	"errors"
	"sync"
	"time"
)

const (
	// Base API v1 Path
	ApiUrlPrefix = "/api/v1/"

	JobPath    = "job/"
	ApiJobPath = ApiUrlPrefix + JobPath

	contentType     = "Content-Type"
	jsonContentType = "application/json;charset=UTF-8"
)

var (
	JobNotFound      = errors.New("Job not found")
	JobCreationError = errors.New("Error creating job")

	GenericError = errors.New("An error occured performing your request")
)

type Duration struct {
	Years   int
	Months  int
	Weeks   int
	Days    int
	Hours   int
	Minutes int
	Seconds int
}

type Job struct {
	Name string `json:"name"`
	Id   string `json:"id"`

	// Command to run
	// e.g. "bash /path/to/my/script.sh"
	Command string `json:"command"`

	// Email of the owner of this job
	// e.g. "admin@example.com"
	Owner string `json:"owner"`

	// Is this job disabled?
	Disabled bool `json:"disabled"`

	// Jobs that are dependent upon this one will be run after this job runs.
	DependentJobs []string `json:"dependent_jobs"`

	// List of ids of jobs that this job is dependent upon.
	ParentJobs []string `json:"parent_jobs"`

	// ISO 8601 String
	// e.g. "R/2014-03-08T20:00:00.000Z/PT2H"
	Schedule     string `json:"schedule"`
	scheduleTime time.Time
	// ISO 8601 Duration struct, used for scheduling
	// job after each run.
	delayDuration *Duration

	// Number of times to schedule this job after the
	// first run.
	timesToRepeat int64

	// Number of times to retry on failed attempt for each run.
	Retries        uint `json:"retries"`
	currentRetries uint

	// Duration in which it is safe to retry the Job.
	Epsilon         string `json:"epsilon"`
	epsilonDuration *Duration

	// Meta data about successful and failed runs.
	SuccessCount     uint      `json:"success_count"`
	LastSuccess      time.Time `json:"last_success"`
	ErrorCount       uint      `json:"error_count"`
	LastError        time.Time `json:"last_error"`
	LastAttemptedRun time.Time `json:"last_attempted_run"`

	jobTimer  *time.Timer
	NextRunAt time.Time `json:"next_run_at"`

	currentStat *JobStat
	Stats       []*JobStat `json:"-"`

	lock sync.Mutex
}

type JobStat struct {
	JobId             string
	RanAt             time.Time
	NumberOfRetries   uint
	Success           bool
	ExecutionDuration time.Duration
}

type JobResponse struct {
	Job *Job `json:"job"`
}

type AddJobResponse struct {
	Id string `json:"id"`
}

type ListJobsResponse struct {
	Jobs map[string]*Job `json:"jobs"`
}

type ListJobStatsResponse struct {
	JobStats []*JobStat `json:"job_stats"`
}

type KalaStats struct {
	ActiveJobs   int
	DisabledJobs int
	Jobs         int

	ErrorCount   uint
	SuccessCount uint

	NextRunAt        time.Time
	LastAttemptedRun time.Time

	CreatedAt time.Time
}

type KalaStatsResponse struct {
	Stats *KalaStats
}
