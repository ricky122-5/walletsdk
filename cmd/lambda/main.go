package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/rickyreddygari/walletsdk/internal/app"
)

var handler http.Handler

func init() {
	container, err := app.NewContainer()
	if err != nil {
		log.Fatalf("bootstrap container: %v", err)
	}
	handler = container.HTTPServer
}

func main() {
	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		w := newLambdaResponseWriter()
		r, err := http.NewRequestWithContext(ctx, request.HTTPMethod, request.Path, strings.NewReader(request.Body))
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
		}

		for key, value := range request.Headers {
			r.Header.Set(key, value)
		}

		if request.QueryStringParameters != nil {
			q := r.URL.Query()
			for key, value := range request.QueryStringParameters {
				q.Set(key, value)
			}
			r.URL.RawQuery = q.Encode()
		}

		handler.ServeHTTP(w, r)

        return events.APIGatewayProxyResponse{
            StatusCode: w.status,
            Headers:    headerMap(w.headers),
            Body:       w.body.String(),
        }, nil
	})
}

type responseWriter struct {
    headers http.Header
	body    strings.Builder
	status  int
}

func (w *responseWriter) Header() http.Header {
    if w.headers == nil {
        w.headers = make(http.Header)
    }
    return w.headers
}

func (w *responseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

func newLambdaResponseWriter() *responseWriter {
    return &responseWriter{headers: make(http.Header), status: http.StatusOK}
}

func headerMap(h http.Header) map[string]string {
    out := make(map[string]string, len(h))
    for k, v := range h {
        if len(v) > 0 {
            out[k] = v[0]
        }
    }
    return out
}
