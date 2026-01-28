package responses

type JobDetails struct {
	WorkerID           string  `json:"worker_id"`
	WorkerStatus       string  `json:"worker_status"`
	WorkerCategories   string  `json:"worker_categories"`
	JobID              string  `json:"job_id"`
	JobStatus          string  `json:"job_status"`
	JobLatitude        float64 `json:"job_latitude"`
	JobLongitude       float64 `json:"job_longitude"`
	JobStartTime       string  `json:"job_start_time"`
	JobDurationHours   float64 `json:"job_duration_hours"`
	JobCategories      string  `json:"job_categories"`
	MatchingCategories string  `json:"matching_categories"`
	ContractorID       string  `json:"contractor_id"`
	ContractorCompany  string  `json:"contractor_company"`
	ProcessID          string  `json:"process_id"`
	ProcessStatus      string  `json:"process_status"`
	DistanceMeters     float64 `json:"distance_meters"`
}

type GetJobsByWorkerResponse struct {
	Jobs []JobDetails `json:"jobs"`
}
