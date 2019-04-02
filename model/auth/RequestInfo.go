package auth

import (
	"net/http"
	// "fmt"
)

type RequestInfo struct{
	Authorization string
	Accept		  string 
	UserAgent	  string
	Method		  string
}


func (r *RequestInfo) SetHeaders(request *http.Request){

	r.Authorization = ""
	r.Accept = ""
	r.UserAgent = ""
	
	if request.Header["Authorization"] != nil {
		r.Authorization = request.Header["Authorization"][0]
	}

	if request.Header["Accept"] != nil {
		r.Accept = request.Header["Accept"][0]
	}

	if request.Header["User-Agent"] != nil {
		r.UserAgent = request.Header["User-Agent"][0]
	}
}

func (r *RequestInfo) GetHttpMethodRequested(request *http.Request) string{
	return request.Method
}