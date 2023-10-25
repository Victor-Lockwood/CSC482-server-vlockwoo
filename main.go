package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func main() {
	gorillamuxRouter := mux.NewRouter()
	gorillamuxRouter.HandleFunc("/vlockwoo/status", status).Methods(http.MethodGet)

	//h := http.HandlerFunc(notFound)
	gorillamuxRouter.NotFoundHandler = gorillamuxRouter.NewRoute().HandlerFunc(http.NotFound).GetHandler()

	gorillamuxRouter.Use(Middleware(gorillamuxRouter))

	http.ListenAndServe(":8090", gorillamuxRouter)
}

func status(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("hello this is me flying gorilla\n"))
	return
}

func notFound(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Not found\n"))
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

func Middleware(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			sw := NewStatusResponseWriter(w)

			defer func() {
				log.Printf(
					"[%s] [%v] [%d] %s %s",
					req.Method,
					time.Since(start),
					sw.statusCode,
					req.RemoteAddr,
					req.URL.Path,
				)
			}()

			next.ServeHTTP(w, req)
		})
	}
}
