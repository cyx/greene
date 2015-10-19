package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cyx/greene"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)

	server := &http.Server{
		Addr:      ":8000",
		Handler:   mux,
		ConnState: greene.New(time.Second * 5),
	}
	server.ListenAndServe()
}

func Home(w http.ResponseWriter, r *http.Request) {
	if _, ok := w.(http.Flusher); !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	for i := 0; i <= 1000; i++ {
		w.Write([]byte(fmt.Sprintf("%d...\n\n", i)))
		w.(http.Flusher).Flush()

		time.Sleep(time.Millisecond * 1)
	}
}
