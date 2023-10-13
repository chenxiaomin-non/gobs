package gethtml

import (
	"io"
	"net/http"
)

type HttpRequest struct {
	Url    string            `json:"url"`
	Method string            `json:"method"`
	Header map[string]string `json:"header"`
	Body   string            `json:"body"`
}

type HttpResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type BodyContentReader struct {
	Body string
}

func (b BodyContentReader) Read(p []byte) (n int, err error) {
	return len(b.Body), io.EOF
}

func GetHttp(request HttpRequest) (HttpResponse, error) {
	httpRequest, err := http.NewRequest(
		request.Method,
		request.Url,
		BodyContentReader{Body: request.Body},
	)

	if err != nil {
		return HttpResponse{}, err
	}

	for key, value := range request.Header {
		httpRequest.Header.Add(key, value)
	}

	client := &http.Client{}
	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		return HttpResponse{}, err
	}

	defer httpResponse.Body.Close()

	var body []byte
	_, err = httpResponse.Body.Read(body)
	if err != nil {
		return HttpResponse{}, err
	}

	return HttpResponse{
		StatusCode: httpResponse.StatusCode,
		Body:       string(body),
	}, nil
}
