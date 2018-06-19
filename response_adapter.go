package celerity

import (
	"encoding/json"
)

// ResponseAdapter - Response adapters are used to marshal an endpoint response
// into bytes
type ResponseAdapter interface {
	Process(Context, Response) ([]byte, error)
}

// JSONResponseAdapter - Processes an endpoint response into JSON
type JSONResponseAdapter struct{}

//JSONResponse - used by the JSONResponseAdapter to build the structure of the
//default JSON response.
type JSONResponse struct {
	RequestID string                 `json:"requestId"`
	Success   bool                   `json:"success"`
	Error     string                 `json:"error"`
	Data      interface{}            `json:"data"`
	Meta      map[string]interface{} `json:"meta"`
}

//Process - Process the response into JSON data.
func (ra *JSONResponseAdapter) Process(c Context, r Response) ([]byte, error) {
	rObj := JSONResponse{
		RequestID: c.RequestID,
		Meta:      r.Meta,
		Data:      r.Data,
	}
	if r.Error != nil {
		rObj.Success = false
		rObj.Error = r.Error.Error()
	} else {
		rObj.Success = true
	}
	return json.Marshal(rObj)
}
