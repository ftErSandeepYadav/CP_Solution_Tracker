package platform

import "time"

type Submission struct{
	ID string
	ProblemID string
	ProblemName string
	ContestID string
	Verdict string
	Language string
	Code string
	SubmissionTime time.Time
	MemoryUsedBytes int64
	TimeConsumedMillis int64
}

type Problem struct{
	ID string
	Name string
	ContestID string
	Index string
	Tags []string
	Difficulty int
	URL string
}

type Platform interface {
	GetName() string
	GetAcceptedSubmissions(handle string)([]Submission, error)
	GetProblemMetadata(problemID string)(*Problem, error)
	GetSubmissionCode(submissionID string)(string, error)
}

type Config struct{
	Handle string
	Cookies string
	APIKey string
	Extra map[string]string
}