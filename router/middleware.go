package router

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func logMiddleware(lg *log.Logger, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nw := time.Now()
		handler.ServeHTTP(w, r)
		lg.Printf("%s %s %s %s\n", r.RemoteAddr, r.Method, r.URL, time.Since(nw))
	})
}

func recoverMiddleware(lg *log.Logger, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("RECOVERING: %v\n", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		handler.ServeHTTP(w, r)
	})
}
