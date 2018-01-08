package grc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type (
	api_info struct {
		server *RestServer
		method string
		uri    string
		body   bool
		proc   func(*http.Request)
		client http.Client
		url    string
	}

	json_api_info struct {
		api_info
	}
)

func (info *api_info) split_body(obj []interface{}) (interface{}, []interface{}) {
	if info.body {
		if len(obj) > 0 {
			return obj[0], obj[1:]
		} else {
			return nil, []interface{}{}
		}
	} else {
		return nil, obj
	}
}

func (info *api_info) get_request(body io.Reader, obj ...interface{}) (*http.Request, error) {
	if info.url == "" {
		if strings.HasPrefix(info.uri, "/") {
			info.url = fmt.Sprintf("%s%s", info.server.get_url(), info.uri)
		} else {
			info.url = fmt.Sprintf("%s/%s", info.server.get_url(), info.uri)
		}
	}
	url := fmt.Sprintf(info.url, obj...)
	fmt.Println("onfo.URL = ", info.url)
	fmt.Println("URL = ", url)
	return http.NewRequest(info.method, url, body)
}

func (info *api_info) run(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	if info.proc != nil {
		info.proc(req)
	}
	return client.Do(req)
}

func (info *json_api_info) Run(obj ...interface{}) (Response, error) {
	b, a := info.split_body(obj)
	var body io.Reader
	if b != nil {

		d, err := json.Marshal(b)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(d)
	}

	req, err := info.get_request(body, a...)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	if info.server.ungzip {
		req.Header.Set("Accept-Encoding", "gzip, deflate")
	}
	fmt.Println(req.Header)
	resp, err := info.run(req)
	fmt.Println(resp.Header)
	if err != nil {
		return nil, err
	}
	return new_grc_response(resp, true), nil
}

func new_json_api(rs *RestServer, method, uri string, body bool, f ...func(*http.Request)) API {
	info := &json_api_info{}
	info.server = rs
	info.method = method
	info.uri = uri
	info.body = body
	if len(f) > 0 {
		info.proc = f[0]
	}
	return info
}
