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
	// The API version
	//
	// required: true
	// example: 0.0.1
	APIVersion string `json:"apiVersion"`
	// The request ID
	//
	// required: true
	// example: dc380b72-41c9-47bf-8be5-f3a7a493f4ca
	ID string `json:"id,omitempty"`
	// The request method
	//
	// required: true
	Method string `json:"method,omitempty"`
}

// Definition is the structure of every http response
type Definition struct {
	Base
	Payload
}
