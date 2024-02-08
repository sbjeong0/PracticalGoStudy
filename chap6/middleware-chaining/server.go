package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type appConfig struct {
	logger *log.Logger
}

type app struct {
	config  appConfig
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		config appConfig,
	)
}

func (a app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler(w, r, a.config)
}

func apiHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	config.logger.Println("Handling API requests")
	fmt.Fprintf(w, "Hello, world\n")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	config.logger.Println("Handling healthcheck request")
	fmt.Fprintf(w, "ok\n")
}

func panicHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	panic("panic!")
}

func setupHandlers(mux *http.ServeMux, config appConfig) {
	mux.Handle("/healthz", &app{config, healthCheckHandler})
	mux.Handle("/api", &app{config, apiHandler})
	mux.Handle("/panic", &app{config, panicHandler})
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			h.ServeHTTP(w, r)
			log.Printf(
				"path=%s method=%s duration=%f",
				r.URL.Path, r.Method, time.Now().Sub(startTime).Seconds(),
			)
		})
}

func panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rValue := recover(); rValue != nil {
					log.Println("panic detected", rValue)
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Unexpected server error")
				}
			}()
			next.ServeHTTP(w, r)
		})
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	config := appConfig{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
	}

	mux := http.NewServeMux()
	setupHandlers(mux, config)

	log.Fatal(http.ListenAndServe(listenAddr, loggingMiddleware(panicMiddleware(mux))))
}
