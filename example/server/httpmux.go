package main

import (
	"io"
	"net/http"
)

//Define a map to implement routing table.
var mux map[string]func(http.ResponseWriter, *http.Request)

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Implement route forwarding
	if h, ok := mux[r.URL.String()]; ok {
		//Implement route forwarding with this handler, the corresponding route calls the corresponding func.
		h(w, r)
		return
	}
	_, _ = io.WriteString(w, "unknown URL: "+r.URL.String())
}
