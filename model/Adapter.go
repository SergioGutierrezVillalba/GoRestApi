package model

import (
	// "net/http"
)

type Adapter struct {}

// func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
// 	for _, adapter := range adapters {
// 	  h = adapter(h)
// 	}
// 	return h
// }

// func (adapter *Adapter) Adapt(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
// 		h.ServeHTTP(w, r)
// 	})
// }