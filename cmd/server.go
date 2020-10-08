package cmd

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)


type Server struct {
	sync.Mutex
	Logging      *Logging
	RequestCount int
	Counters     map[string]int
}

func (server Server) start(port string) {
	http.HandleFunc("/counter", server.counter)
	http.HandleFunc("/counter-without-mutex", server.counterWithoutMutex)
	http.HandleFunc("/counter-without-pointer", server.counterWithoutPointer)
	server.Logging.Stdout.Printf("Start server, listing on port: %s", port)

	server.Counters["requests"] = 0

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		server.Logging.ErrorLog.Printf("%v", err)
	}
}

func (server *Server) counter(w http.ResponseWriter, req *http.Request) {
	time.Sleep(3 * time.Millisecond)
	server.Lock()
	defer server.Unlock()
	server.Counters["requests"]++
	server.Logging.Stdout.Printf("incoming request count %d", server.Counters["requests"])
	_, _ = fmt.Fprintf(w, "ok")
	//http.Error(w, "Unknown Error", http.StatusInternalServerError)
}

func (server *Server) counterWithoutMutex(w http.ResponseWriter, req *http.Request) {
	time.Sleep(3 * time.Millisecond)
	server.Counters["requests"]++
	server.Logging.Stdout.Printf("incoming request count %d", server.Counters["requests"])
	_, _ = fmt.Fprintf(w, "ok")
}

func (server Server) counterWithoutPointer(w http.ResponseWriter, req *http.Request) {
	time.Sleep(3 * time.Millisecond)
	server.Lock()
	defer server.Unlock()
	server.Counters["requests"]++
	server.Logging.Stdout.Printf("incoming request count %d", server.Counters["requests"])
	_, _ = fmt.Fprintf(w, "ok")
}
