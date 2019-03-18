package auth

import (
	"net/http"
	// "fmt"
)

type RequestHeaders struct{
	Authorization string
	Accept		  string 
	UserAgent	  string
}


func (requestHeaders *RequestHeaders) SetHeaders(r *http.Request){

	requestHeaders.Authorization = ""
	requestHeaders.Accept = ""
	requestHeaders.UserAgent = ""
	
	if r.Header["Authorization"] != nil {
		requestHeaders.Authorization = r.Header["Authorization"][0]
	}

	if r.Header["Accept"] != nil {
		requestHeaders.Accept = r.Header["Accept"][0]
	}

	if r.Header["User-Agent"] != nil {
		requestHeaders.UserAgent = r.Header["User-Agent"][0]
	}

}