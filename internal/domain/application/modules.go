package application

// Module bindings
const (
	ModuleApp             string = "modules.application"
	ModuleLogger          string = "modules.logger"
	ModuleBaseLogger      string = "modules.base-logger"
	ModuleSchemaRegistry  string = "modules.schema-registry"
	ModuleMetricsReporter string = "modules.metrics-reporter"
	ModuleBaseReporter    string = "modules.base-metrics-reporter"
	ModuleSQL             string = "module-sql"
	MoudleDBConector      string = "module-das-db-connector"

	ModuleHTTP           string = "modules.http"
	ModuleHTTPRouter     string = "modules.router"
	ModuleHTTPServer     string = "modules.http.server"
	ModuleReadyIndicator string = "module-ready-indicator"
	ModulePprofIndicator string = "module-pprof-indicator"
	ModuleErrorHandler   string = "http-error-handler"
)

// Encoder
const (
	ModuleEncoders      string = "encoders"
	ModuleJSONEncoder   string = "encoders.json"
	ModuleStringEncoder string = "encoders.string"
	ModuleUUIDEncoder   string = "encoders.uuid"
)

// Repositories
const (
	ModuleAccountRepo       string = "repositories.account"
	ModulePaymentRepo       string = "repositories.payment"
	ModuleAccountWalletRepo string = "repositories.account-wallet"
)

// Services
const (
	ModulePaymentService string = "services.payment"
	ModuleWalletService  string = "services.wallet"
)
