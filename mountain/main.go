package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"net/http"
)

type Response struct {
	LongestPeak int `json:"longestPeak"`
}

var httpLambda *httpadapter.HandlerAdapter

func convertStringArrayToIntArray(stringArray string) ([]int, error) {
	var intArray []int
	err := json.Unmarshal([]byte(stringArray), &intArray)
	if err != nil {
		return nil, err
	}
	return intArray, nil
}

func mountain(w http.ResponseWriter, r *http.Request) {
	stringArray := r.URL.Query().Get("mountainArray")
	intArray, err := convertStringArrayToIntArray(stringArray)
	if err != nil {
		http.Error(w, "Error converting query string to int array.", 400)
		return
	}

	var response Response
	mountainResult := LongestMountain(intArray)
	response.LongestPeak = mountainResult

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error converting response struct to json.", 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonResponse)
}

func info(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a mountain API."))
}

func init() {
	http.HandleFunc("GET /mountain", mountain)
	http.HandleFunc("GET /info", info)
	httpLambda = httpadapter.New(http.DefaultServeMux)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return httpLambda.ProxyWithContext(ctx, req)
}

func main() {
	http.ListenAndServe("127.0.0.1:8080", nil)
	//lambda.Start(Handler)
}
