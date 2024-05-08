package mdlwr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

func LogMiddleware(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
		logRequest(&ctx, &request, nil, nil)
		addResponseFunc(&ctx)
		return next(ctx, request)
	}
}

type LogMiddlewareFunc func(HandlerFunc, *[]string, *[]string) HandlerFunc

// Returns a new logger middleware that exclude the keys passed as parameters

// To filter a key inside a nested map, you have to concatenate the keys with "."

// For example, to tell the logger to exclude the content key inside the body,
// we have to pass "files.content" inside the array for the param ignoreBodyLog
//
// And if we want to exclude also the "type" we would make an array including both values
//
//	 []string{"files.content", "type"}
//
//	 {
//	     "files": [
//	         {
//				"fileName": "file1",
//				"content": "lkasdjalskds.."
//	         },
//	         {
//				"fileName": "file2",
//				"content": "askmlmalksiod.."
//	         }
//	         "type": "IIBB",
//	     ],
//	 }
//
// For excluding values of the request header, we have to pass also
// an array of strings with the key we want to hide.
//
// For example, if we want to hide the "content-length" key,
// we would pass this array as second argument value
//
//	[]string{"content-length"}
func CreateLogMdlw(ignoreBodyLog *[]string, ignoreHeaderLog *[]string) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
			logRequest(&ctx, &request, ignoreBodyLog, ignoreHeaderLog)
			addResponseFunc(&ctx)
			return next(ctx, request)
		}
	}
}

func logRequest(ctx *context.Context, request *events.APIGatewayV2HTTPRequest, ignoreBodyLog *[]string, ignoreHeaderLog *[]string) {
	body := make(map[string]interface{})
	json.Unmarshal([]byte(request.Body), &body)

	log.Info().Str("trace_id", fmt.Sprint((*ctx).Value("traceId"))).
		Str("method", request.RequestContext.HTTP.Method).
		Str("path", request.RequestContext.HTTP.Path).
		RawJSON("query_params", getQueryParams(request.QueryStringParameters)).
		RawJSON("header_params", getHeaders(request.Headers, ignoreHeaderLog)).
		Msg("Incoming request " + getBody(body, ignoreBodyLog))
}

func getBody(body map[string]interface{}, ignoreBodyLog *[]string) string {
	original := reflect.ValueOf(body)
	copy := reflect.New(original.Type()).Elem()
	if ignoreBodyLog != nil {
		for _, v := range *ignoreBodyLog {
			arrayExclude := strings.Split(v, ".")
			filterRecursive(copy, original, arrayExclude, 0)
			original = copy
			copy = reflect.New(original.Type()).Elem()
		}
	}

	filteredBodyJson, _ := json.Marshal(original.Interface())

	var b bytes.Buffer
	json.HTMLEscape(&b, filteredBodyJson)
	return b.String()
}

func getQueryParams(queryParam map[string]string) []byte {
	jqp, _ := json.Marshal(queryParam)
	return jqp
}

func addResponseFunc(ctx *context.Context) {
	*ctx = context.WithValue(*ctx, "logRepsonseFunc", ResponseFunc)
}

func ResponseFunc(ctx *context.Context, resp *events.APIGatewayProxyResponse) {
	log.Info().
		Str("trace_id", fmt.Sprint((*ctx).Value("traceId"))).
		RawJSON("header", getHeaders(resp.Headers, nil)). //Header params del response
		Msg("Response status code: " + fmt.Sprint(resp.StatusCode))
}

func getHeaders(headers map[string]string, ignoreHeaderLog *[]string) []byte {
	h := make(map[string]string)
	h = filterHeader(headers, ignoreHeaderLog)
	jh, _ := json.Marshal(h)
	return jh
}

func filterHeader(headers map[string]string, ignoreHeaderLog *[]string) map[string]string {
	h := make(map[string]string)
	for param, value := range headers {
		if !strings.HasPrefix(strings.ToLower(param), "authorization") {
			h[param] = value
		}
	}
	if ignoreHeaderLog != nil {
		for _, value := range *ignoreHeaderLog {
			if _, ok := h[value]; ok {
				delete(h, value)
			}

		}
	}

	return h
}

func filterRecursive(copy, original reflect.Value, exclude []string, pos int) {
	if pos >= len(exclude) {
		return
	}
	switch original.Kind() {
	case reflect.Ptr:
		originalValue := original.Elem()
		if !originalValue.IsValid() {
			return
		}
		copy.Set(reflect.New(originalValue.Type()))
		filterRecursive(copy.Elem(), originalValue, exclude, pos)
	case reflect.Interface:
		if original.IsZero() || original.IsNil() {
			break
		}
		originalValue := original.Elem()
		copyValue := reflect.New(originalValue.Type()).Elem()
		filterRecursive(copyValue, originalValue, exclude, pos)
		copy.Set(copyValue)
	case reflect.Struct:
		for i := 0; i < original.NumField(); i++ {
			filterRecursive(copy.Field(i), original.Field(i), exclude, pos)
		}
	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i++ {
			filterRecursive(copy.Index(i), original.Index(i), exclude, pos)
		}
	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			if exclude[pos] == key.String() {
				filterRecursive(copyValue, originalValue, exclude, pos+1)
				if pos != len(exclude)-1 {
					copy.SetMapIndex(key, copyValue)
				}
			} else {
				filterRecursive(copyValue, originalValue, exclude, pos)
				copy.SetMapIndex(key, copyValue)
			}
		}
	default:
		copy.Set(original)
	}
}
