package lambdasqs

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
)

const PointerErrorString = "object param must be a pointer"

func BindData(message events.SQSMessage, object interface{}) error {
	ro := reflect.ValueOf(object)
	if ro.Kind() != reflect.Ptr || ro.IsNil() {
		return errors.New(PointerErrorString)
	}
	err := ExtractBody(message, object)
	if err != nil {
		return err
	}
	return nil
}

func ExtractBody(message events.SQSMessage, object interface{}) error {
	err := json.Unmarshal([]byte(message.Body), object)
	if err != nil {
		return err
	}
	return nil
}
