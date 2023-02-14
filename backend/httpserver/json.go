package httpserver

import (
	"encoding/json"
	"log"
)

// JsonResponse
type JsonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// SetJsonResponse
func (r *RequestContext) SetJsonResponse(res interface{}) error {
	r.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(r.ResponseWriter)
	encoder.SetEscapeHTML(true)
	log.Printf("%v %v %v %+v\n", r.Request.RemoteAddr, r.Request.Method, r.Request.URL, res)
	return encoder.Encode(res)
}
