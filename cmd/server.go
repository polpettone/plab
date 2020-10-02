package cmd

import (
	"fmt"
	"net/http"
)

func (server Server) ok(w http.ResponseWriter, req *http.Request) {
	server.Logging.stdout.Printf(req.RemoteAddr)
	_, _ = fmt.Fprintf(w, "ok")
}

type Server struct {
	Logging *Logging
}

func (server Server) start() {
	http.HandleFunc("/", server.ok)
	http.ListenAndServe(":8080", nil)
}