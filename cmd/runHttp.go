package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chilledoj/mystdhttp/router"
)

type options struct {
	listenAddr string
	initTasks  bool
}

func runHttp(opts options) error {
	s := http.Server{
		Addr:           opts.listenAddr,
		Handler:        router.NewRouter(opts.initTasks),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("Starting HTTP listener at %s\n", opts.listenAddr)
	return s.ListenAndServe()
}
