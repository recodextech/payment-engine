package domain

const (
	AccountTable                  string = `"fix.account"`
	ContractorTable               string = `"fix.contractor"`
	WorkerTable                   string = `"fix.worker"`
	ProcessTable                  string = `"fix.process"`
	JobTable                      string = `"fix.job"`
	JobLocationTable              string = `"fix.job.location"`
	JobAssignmentTable            string = `"fix.job.assignment"`
	ConfigurationCategory         string = `"fix.configuration.category"`
	JobCategoryTable              string = `"fix.job.category"`
	WorkerCategoryTable           string = `"fix.worker.category"`
	WorkerAvailabilityTable       string = `"fix.worker.availability"`
	WorkerAvailabilityWindowTable string = `"fix.worker.availability.window"`
	
	// Views
	WorkerNearbyJobsView          string = `"vw_worker_nearby_jobs"`
)
