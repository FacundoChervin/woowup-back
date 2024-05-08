package lambdahttp

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mitchellh/mapstructure"
)

func BindData(request events.APIGatewayV2HTTPRequest, object interface{}) error {
	ExtractBody(request, object)
	ExtractParams(request, object)
	ExtractHeaders(request, object)
	ExtractPathParams(request, object)

	return nil
}

func ExtractBody(request events.APIGatewayV2HTTPRequest, object interface{}) error {
	err := json.Unmarshal([]byte(request.Body), &object)
	if err != nil {
		return err
	}
	return nil
}

func ExtractParams(request events.APIGatewayV2HTTPRequest, object interface{}) error {
	decoderConfig := mapstructure.DecoderConfig{
		WeaklyTypedInput:     true,
		TagName:              "query",
		Result:               &object,
		IgnoreUntaggedFields: true,
	}

	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		return err
	}

	if request.QueryStringParameters != nil {
		err = decoder.Decode(request.QueryStringParameters)
	}
	if err != nil {
		return err
	}

	return nil
}

func ExtractHeaders(request events.APIGatewayV2HTTPRequest, object interface{}) error {
	decoderConfig := mapstructure.DecoderConfig{
		WeaklyTypedInput:     true,
		TagName:              "header",
		Result:               &object,
		IgnoreUntaggedFields: true,
	}

	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		return err
	}

	if request.Headers != nil {
		err = decoder.Decode(request.Headers)
	}

	if err != nil {
		return err
	}
	return nil
}

func ExtractPathParams(request events.APIGatewayV2HTTPRequest, object interface{}) error {
	decoderConfig := mapstructure.DecoderConfig{
		WeaklyTypedInput:     true,
		TagName:              "param",
		Result:               &object,
		IgnoreUntaggedFields: true,
	}

	decoder, err := mapstructure.NewDecoder(&decoderConfig)
	if err != nil {
		return err
	}

	if request.PathParameters != nil {
		err = decoder.Decode(request.PathParameters)
	}

	if err != nil {
		return err
	}
	return nil
}
