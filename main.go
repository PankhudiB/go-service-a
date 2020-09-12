package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-service/constants"
	"go-service/request"
	"go-service/tracing"
	"go.opencensus.io/plugin/ochttp"
	"io/ioutil"
	"net/http"
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

func handleRequest(ctx *gin.Context) error {
	fmt.Println("------------------ Welcome to Service A ------------------")
	var serviceARequest request.HelloARequest
	if err := ctx.ShouldBindJSON(&serviceARequest); err != nil {
		fmt.Println("Error reading request hello A", err.Error())
	}
	fmt.Printf("%s said Hello", serviceARequest.Sender)
	err := callB(ctx, serviceARequest)
	if err != nil {
		return err
	}
	return nil
}

func callB(ctx *gin.Context, serviceARequest request.HelloARequest) error {
	bRequest := request.NewHelloBRequest("Sender A", fmt.Sprintf("Just came by to say hii on behalf of %s!", serviceARequest.Sender))
	httpRequest, err := bRequest.CreateHTTPRequest(ctx)

	client := &http.Client{Transport: &ochttp.Transport{}}
	response, err := client.Do(httpRequest)
	if err != nil {
		fmt.Println("Error while requesting service B: ", err.Error())
		return err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error while reading response from service B: ", err.Error())
		return err
	}
	fmt.Println(string(body))
	return nil
}
