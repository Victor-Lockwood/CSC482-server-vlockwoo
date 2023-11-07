package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/gorilla/mux"
	"github.com/jamespearly/loggly"
	"net/http"
	"strconv"
	"time"
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

func all(w http.ResponseWriter, req *http.Request) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	tableName := "vlockwoo-satellites"

	sess.Config.Region = aws.String("us-east-1")
	// Create DynamoDB client
	svc := dynamodb.New(sess)

	input := dynamodb.ScanInput{TableName: aws.String(tableName)}
	output, _ := svc.Scan(&input)

	finalOutput := output.Items

	//TODO: I don't have enough data to test this and I'm not comfortable pushing this uncommented until I do.  Currently there are 54 entries in the DB and I don't want to push my limits just yet.
	//for {
	//	if output.LastEvaluatedKey != nil {
	//		input := dynamodb.ScanInput{TableName: aws.String(tableName), ExclusiveStartKey: output.LastEvaluatedKey}
	//		output, _ := svc.Scan(&input)
	//
	//		finalOutput = append(finalOutput, output.Items...)
	//	} else {
	//		break
	//	}
	//}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(finalOutput)
	return
}

// Search for items based on either:
// eclipsed 	- Whether the space station is currently eclipsed relative to Rice Creek
// timestamp 	- A known timestamp
// Most other attributes are very specific (eg. azimuth) or currently unchanging (eg. satid).
func search(w http.ResponseWriter, req *http.Request) {
	//puts query list into a map
	query := req.URL.Query()

	//check number of parameters
	queryLen := len(query)

	if queryLen == 0 {
		endpointResponse := EndpointResponse{SystemTime: time.Now(), Status: http.StatusBadRequest}

		start := time.Now()
		sw := NewStatusResponseWriter(w)
		sw.WriteHeader(http.StatusBadRequest)

		logRequestToLoggly(sw, req, start)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(endpointResponse)
		return
	}

	eclipsed, present := query["eclipsed"]

	if present && len(eclipsed[0]) == 0 {
		fmt.Println("\teclipsed specified but not present\n")

		endpointResponse := EndpointResponse{SystemTime: time.Now(), Status: http.StatusBadRequest}

		start := time.Now()
		sw := NewStatusResponseWriter(w)
		sw.WriteHeader(http.StatusBadRequest)

		logRequestToLoggly(sw, req, start)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(endpointResponse)
		return
	} else if present {
		isEclipsed, err := strconv.ParseBool(eclipsed[0])

		if err != nil {
			fmt.Println("\teclipsed is not a bool - 400")

			endpointResponse := new(EndpointResponse)
			endpointResponse.SystemTime = time.Now()
			endpointResponse.Status = http.StatusBadRequest

			start := time.Now()
			sw := NewStatusResponseWriter(w)
			sw.WriteHeader(http.StatusBadRequest)

			logRequestToLoggly(sw, req, start)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(endpointResponse)
			return
		} else {
			sess := session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))

			tableName := "vlockwoo-satellites"

			sess.Config.Region = aws.String("us-east-1")
			// Create DynamoDB client
			svc := dynamodb.New(sess)

			filt := expression.Name("eclipsed").Equal(expression.Value(isEclipsed))
			expr, _ := expression.NewBuilder().WithFilter(filt).Build()

			params := &dynamodb.ScanInput{
				TableName:                 aws.String(tableName),
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				FilterExpression:          expr.Filter(),
			}

			output, _ := svc.Scan(params)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(output.Items)
			return
		}
	}

	//Searching by Timestamp
	timestamp, present := query["timestamp"]

	if present && len(timestamp[0]) == 0 {
		endpointResponse := EndpointResponse{SystemTime: time.Now(), Status: http.StatusBadRequest}

		start := time.Now()
		sw := NewStatusResponseWriter(w)
		sw.WriteHeader(http.StatusBadRequest)

		logRequestToLoggly(sw, req, start)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(endpointResponse)
		return
	} else if present {
		specifiedTimestamp, err := strconv.Atoi(timestamp[0])

		if err != nil || !(specifiedTimestamp > 0) {
			fmt.Println("\ttimestamp is not an int/not greater than 0 - 400")

			endpointResponse := new(EndpointResponse)
			endpointResponse.SystemTime = time.Now()
			endpointResponse.Status = http.StatusBadRequest

			start := time.Now()
			sw := NewStatusResponseWriter(w)
			sw.WriteHeader(http.StatusBadRequest)

			logRequestToLoggly(sw, req, start)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(endpointResponse)
			return
		} else {
			sess := session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))

			tableName := "vlockwoo-satellites"

			sess.Config.Region = aws.String("us-east-1")
			// Create DynamoDB client
			svc := dynamodb.New(sess)

			filt := expression.Name("timestamp").Equal(expression.Value(specifiedTimestamp))
			expr, _ := expression.NewBuilder().WithFilter(filt).Build()

			params := &dynamodb.ScanInput{
				TableName:                 aws.String(tableName),
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				FilterExpression:          expr.Filter(),
			}

			output, _ := svc.Scan(params)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(output.Items)
			return
		}
	}
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
