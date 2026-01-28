package handlers

type HandlerError struct {
	error
}

type ParameterExtractionError struct {
	error
}

type DependencyError struct {
	error
}
