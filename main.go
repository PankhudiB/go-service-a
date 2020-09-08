package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-service/constants"
	"go-service/request"
	"go-service/tracing"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {

	fmt.Print("Starting service A")
	tracing.Init("Service A", constants.OCCollector)
	r := gin.Default()
	r.POST("/hello-service-A", func(ctx *gin.Context) {
		handleRequest(ctx)
	})

	err := http.ListenAndServe(":8080", tracing.WithTracing(r))
	if err != nil {
		fmt.Println("Could not start service A", err)
	}

}

func handleRequest(ctx *gin.Context) {
	fmt.Println("------------------ Welcome to Service A ------------------")
	var serviceARequest request.HelloARequest
	if err := ctx.ShouldBindJSON(&serviceARequest); err != nil {
		fmt.Println("Error reading request hello A", err.Error())
	}
	fmt.Printf("%s said Hello", serviceARequest.Sender)
	callB(ctx, serviceARequest)
}

func callB(ctx *gin.Context, serviceARequest request.HelloARequest) {
	req, err := createRequestForB(ctx, serviceARequest)
	client := &http.Client{Transport: &ochttp.Transport{}}
	response, err := client.Do(req)
	if err != nil {
		print(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		print(err)
	}
	fmt.Println(string(body))
}

func createRequestForB(ctx *gin.Context, serviceARequest request.HelloARequest) (*http.Request, error) {
	context := ctx.Request.Context()
	span := trace.FromContext(context)
	defer span.End()
	span.Annotate([]trace.Attribute{trace.StringAttribute("Creating request", "for service B")}, "Logs")
	time.Sleep(time.Millisecond * 1000)

	requestToB := request.HelloBRequest{
		Sender:  "Service-A",
		Message: fmt.Sprintf("Just came by to say hii on behalf of %s!", serviceARequest.Sender),
	}
	requestBody, err := json.Marshal(requestToB)
	if err != nil {
		fmt.Println("Unable to marshal request ", err.Error())
	}
	req, _ := http.NewRequest("POST", constants.HelloServiceBUrl, bytes.NewBuffer(requestBody))
	req = req.WithContext(context)
	return req, err
}
