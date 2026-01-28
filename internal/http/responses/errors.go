package responses

type ErrorResponse struct {
	Code        int    `json:"code"`
	Status      int    `json:"-"`
	Description string `json:"description"`
	Trace       string `json:"trace"`
	Debug       string `json:"debug,omitempty"`
}
