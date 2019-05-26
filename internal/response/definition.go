package response

// Payload is the payload for every response
type Payload struct {
	Data  []interface{}          `json:"data"`
	Error map[string]interface{} `json:"error"`
}

// Properties are the predefined set of properties for each response
type Properties struct {
	APIVersion string
}

// Definition is the structure of every http response
type Definition struct {
	APIVersion string                 `json:"apiVersion"`
	ID         string                 `json:"id,omitempty"`
	Method     string                 `json:"method,omitempty"`
	Data       []interface{}          `json:"data,omitempty"`
	Errors     map[string]interface{} `json:"error,omitempty"`
}

