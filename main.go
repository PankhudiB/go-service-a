package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-service/request"
	"net/http"
)

func main() {

	fmt.Print("Starting service A")

	r := gin.Default()
	r.POST("/hello-service-A", func(ctx *gin.Context) {
		handleRequest(ctx)
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Could not start service A", err)
	}

}

func handleRequest(ctx *gin.Context) {
	fmt.Println("------------------ Welcome to Service A ------------------")
	var req request.HelloARequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println("Error reading request hello A", err.Error())
	}
	fmt.Printf("-------- %s said Hello --------", req.Sender)
}
