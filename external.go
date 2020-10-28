package restqlnewrelic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
)

func makeExternalRequest(request restql.HTTPRequest) (*http.Request, error) {
	data, err := json.Marshal(request.Body)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s://%s%s", request.Schema, request.Host, request.Path)
	r, err := http.NewRequest(request.Method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	hdr := make(http.Header)
	for k, v := range request.Headers {
		hdr.Add(k, v)
	}
	r.Header = hdr

	return r, nil
}

func makeExternalResponse(request *http.Request, response restql.HTTPResponse) (*http.Response, error) {
	body, err := json.Marshal(request.Body)
	if err != nil {
		return nil, err
	}

	r := http.Response{}
	r.StatusCode = response.StatusCode
	hdr := make(http.Header)
	for k, v := range response.Headers {
		hdr.Add(k, v)
	}
	r.Header = hdr
	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	r.Request = request

	return &r, nil
}
