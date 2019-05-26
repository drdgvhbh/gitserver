package response

// Payload is the payload for every response
type Payload struct {
	Data   []interface{}          `json:"data,omitempty"`
	Errors map[string]interface{} `json:"errors,omitempty"`
}

// Properties are the predefined set of properties for each response
type Properties struct {
	APIVersion string
}

// Base is the basic definition of every response
type Base struct {
	APIVersion string                 `json:"apiVersion"`
	ID         string                 `json:"id,omitempty"`
	Method     string                 `json:"method,omitempty"`
}

// Definition is the structure of every http response
type Definition struct {
	Base
	Payload
}

