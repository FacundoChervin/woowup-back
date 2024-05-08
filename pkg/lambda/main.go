package lambda

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"main/pkg/logger"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var TraceIdNotFoundErr = errors.New("DEFCON 1: no trace id found in context")
var TraceIdByteMessage interface{} = "traceId not found or implemented"

func buildResponse(ctx context.Context, dataMap map[string]interface{}, code int) *events.APIGatewayProxyResponse {
	dataByte, _ := json.Marshal(dataMap)

	var dataBuf bytes.Buffer
	json.HTMLEscape(&dataBuf, dataByte)

	traceId := ctx.Value("traceId")
	if traceId == nil {
		traceId = TraceIdByteMessage
		logger.BackgroundError(TraceIdNotFoundErr)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       dataBuf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Request-Id": traceId.(string),
		},
	}
}

func LambdaErrorResponse(ctx context.Context, err error, code int) *events.APIGatewayProxyResponse {

	errMap := map[string]interface{}{
		"status":  parseStatus(code),
		"message": err.Error(),
	}
	resp := buildResponse(ctx, errMap, code)

	logResponse(&ctx, resp)
	return resp
}

func LambdaSuccessResponse(ctx context.Context, statusCode int, datas ...ResponseData) *events.APIGatewayProxyResponse {
	successMap := map[string]interface{}{
		"status": parseStatus(statusCode),
		"data":   toDataMap(datas),
	}
	resp := buildResponse(ctx, successMap, statusCode)

	logResponse(&ctx, resp)
	return resp
}

func LambdaCreatedResponse(ctx context.Context) *events.APIGatewayProxyResponse {
	traceId := ctx.Value("traceId")
	if traceId == nil {
		traceId = TraceIdByteMessage
		logger.BackgroundError(TraceIdNotFoundErr)
	}
	resp := &events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers: map[string]string{
			"X-Request-Id": traceId.(string),
		},
	}

	logResponse(&ctx, resp)
	return resp
}

func LambdaNoContentResponse(ctx context.Context) *events.APIGatewayProxyResponse {
	traceId := ctx.Value("traceId")
	if traceId == nil {
		traceId = TraceIdByteMessage
		logger.BackgroundError(TraceIdNotFoundErr)
	}
	resp := &events.APIGatewayProxyResponse{
		StatusCode: http.StatusNoContent,
		Headers: map[string]string{
			"X-Request-Id": traceId.(string),
		},
	}

	logResponse(&ctx, resp)
	return resp
}

func logResponse(ctx *context.Context, resp *events.APIGatewayProxyResponse) {
	afterFunc := (*ctx).Value("logRepsonseFunc")
	if afterFunc == nil {
		return
	}
	afterFunc.(func(*context.Context, *events.APIGatewayProxyResponse))(ctx, resp)
}

type ResponseData struct {
	Name string
	Data interface{}
}

func toDataMap(datas []ResponseData) map[string]interface{} {
	if len(datas) < 1 {
		return nil
	}
	dataMap := make(map[string]interface{})
	for _, data := range datas {
		dataMap[data.Name] = data.Data
	}
	return dataMap
}

func parseStatus(code int) string {
	switch code / 100 {
	case 2:
		return "success"
	case 4:
		return "fail"
	case 5:
		return "error"
	}
	return "unknown"
}
