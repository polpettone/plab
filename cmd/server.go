package cmd

import (
	"net/http"
)

func (server Server) ok(w http.ResponseWriter, req *http.Request) {
	server.Logging.stdout.Printf(req.RemoteAddr)
	//_, _ = fmt.Fprintf(w, "ok")
	http.Error(w, "Unknown Error", http.StatusInternalServerError)
}

type Server struct {
	Logging *Logging
}

func (server Server) start(port string) {
	http.HandleFunc("/", server.ok)
	server.Logging.stdout.Printf("Start server, listing on port: %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		server.Logging.errorLog.Printf("%v", err)
	}
}