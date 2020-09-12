package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-service/constants"
	"go.opencensus.io/trace"
	"net/http"
	"time"
)

type HelloBRequest struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

func NewHelloBRequest(sender, message string) HelloBRequest {
	return HelloBRequest{
		Sender:  sender,
		Message: message,
	}
}

func (requestB HelloBRequest) CreateHTTPRequest(ctx *gin.Context) (*http.Request, error) {
	context := ctx.Request.Context()
	span := trace.FromContext(context)
	defer span.End()
	span.Annotate([]trace.Attribute{trace.StringAttribute("Creating request", "for service B")}, "Logs")
	time.Sleep(time.Millisecond * 1000)

	requestBody, err := json.Marshal(requestB)
	if err != nil {
		fmt.Println("Unable to marshal request ", err.Error())
		return nil, err
	}
	req, _ := http.NewRequest("POST", constants.HelloServiceBUrl, bytes.NewBuffer(requestBody))
	req = req.WithContext(context)
	return req, err
}
