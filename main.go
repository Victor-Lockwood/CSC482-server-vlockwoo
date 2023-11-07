package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jamespearly/loggly"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	gorillamuxRouter := mux.NewRouter()

	gorillamuxRouter.HandleFunc("/vlockwoo/status", status).Methods(http.MethodGet)
	gorillamuxRouter.HandleFunc("/vlockwoo/all", all).Methods(http.MethodGet)
	gorillamuxRouter.HandleFunc("/vlockwoo/search", search).Methods(http.MethodGet)

	gorillamuxRouter.NotFoundHandler = http.HandlerFunc(notFound)
	gorillamuxRouter.MethodNotAllowedHandler = http.HandlerFunc(notAllowed)

	gorillamuxRouter.Use(Middleware(gorillamuxRouter))

	http.ListenAndServe(":8080", gorillamuxRouter)
}

// *** ENDPOINTS ***

func status(w http.ResponseWriter, req *http.Request) {
	endpointResponse := new(EndpointResponse)
	endpointResponse.SystemTime = time.Now()
	endpointResponse.Status = http.StatusOK

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(endpointResponse)
	return
}

func all(w http.ResponseWriter, req *http.Request) {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	tableName := "vlockwoo-satellites"

	sess.Config.Region = aws.String("us-east-1")
	// Create DynamoDB client
	svc := dynamodb.New(sess)

	input := dynamodb.ScanInput{TableName: aws.String(tableName)}
	output, err := svc.Scan(&input)

	if err != nil {
		fmt.Printf("Got error")
	}

	response := AllResponse{
		TableName:   tableName,
		RecordCount: *output.Count,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return
}

func search(w http.ResponseWriter, req *http.Request) {
	//puts query list into a map
	query := req.URL.Query()

	//check number of parameters
	queryLen := len(query)
	fmt.Printf("\tLENGTH:%d\n", queryLen)

	//Check parameter names
	fmt.Printf("\t%+v\n", query)

	test1, present := query["test1"]
	fmt.Printf("\ttest1 len %d\n", len(test1))

	if !present || len(test1[0]) == 0 {
		fmt.Println("\ttest1 not present\n")
		//You can return a 400 right here
	}

	//test1 should be an int
	test1Val, err := strconv.Atoi(test1[0])

	if err != nil {
		fmt.Println("\ttest1 is not an int - 400")
	} else {
		fmt.Printf("\ttest1 vale us %d\n", test1Val)
	}

	// Check that all expected vals make sense

	//Check for only alpha numeric
	//re := regexp.MustCompile("^[a-zA-Z]+$")
	//fmt.Printf("\ttest2 REGEX TEST: %v\n", re.MatchString(test2[0]))
}

// *** 404 AND 405 HANDLERS ***

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

// *** MIDDLEWARE ***

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
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
