package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	loggly "github.com/jamespearly/loggly"
	"net/http"
	"time"
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

type EndpointResponse struct {
	SystemTime time.Time `json:"systemtime"`
	Status     int       `json:"status"`
}

func main() {
	gorillamuxRouter := mux.NewRouter()

	gorillamuxRouter.HandleFunc("/vlockwoo/status", status).Methods(http.MethodGet)
	gorillamuxRouter.NotFoundHandler = http.HandlerFunc(notFound)
	gorillamuxRouter.MethodNotAllowedHandler = http.HandlerFunc(notAllowed)

	gorillamuxRouter.Use(Middleware(gorillamuxRouter))

	http.ListenAndServe(":8090", gorillamuxRouter)
}

func status(w http.ResponseWriter, req *http.Request) {
	endpointResponse := new(EndpointResponse)
	endpointResponse.SystemTime = time.Now()
	endpointResponse.Status = http.StatusOK

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(endpointResponse)
	return
}

func notFound(w http.ResponseWriter, req *http.Request) {
	endpointResponse := new(EndpointResponse)
	endpointResponse.SystemTime = time.Now()
	endpointResponse.Status = http.StatusNotFound

	start := time.Now()
	sw := NewStatusResponseWriter(w)
	sw.WriteHeader(http.StatusNotFound)
	logRequestToLoggly(sw, req, start)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(endpointResponse)
	return
}

func notAllowed(w http.ResponseWriter, req *http.Request) {
	endpointResponse := new(EndpointResponse)
	endpointResponse.SystemTime = time.Now()
	endpointResponse.Status = http.StatusMethodNotAllowed

	start := time.Now()
	sw := NewStatusResponseWriter(w)
	sw.WriteHeader(http.StatusMethodNotAllowed)

	logRequestToLoggly(sw, req, start)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(endpointResponse)
	return
}

func NewStatusResponseWriter(responseWriter http.ResponseWriter) *statusResponseWriter {
	return &statusResponseWriter{
		ResponseWriter: responseWriter,
		statusCode:     http.StatusOK,
	}
}

func (sw *statusResponseWriter) WriteHeader(statusCode int) {
	sw.statusCode = statusCode
	sw.ResponseWriter.WriteHeader(statusCode)
}

// Middleware TODO: I am unable to figure out how to make this work for error cases.  I have a workaround but it's not ideal.
func Middleware(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			sw := NewStatusResponseWriter(w)
			logRequestToLoggly(sw, req, start)

			next.ServeHTTP(w, req)
		})
	}
}

func logRequestToLoggly(sw *statusResponseWriter, req *http.Request, start time.Time) {
	message := fmt.Sprintf(
		"[%s] [%v] [%d] %s %s",
		req.Method,
		time.Since(start),
		sw.statusCode,
		req.RemoteAddr,
		req.URL.Path,
	)

	logglyClient := loggly.New("Server")
	_ = logglyClient.EchoSend("info", message)
}
