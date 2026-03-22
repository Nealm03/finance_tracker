package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("FOO"))
	})
	s := &http.Server{
		Addr:         ":9090",
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		Handler:      mux,
	}

	if err := s.ListenAndServe(); err != nil {
		fmt.Println(fmt.Errorf("server err: %w", err))
	}
}
