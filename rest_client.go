package grc

import (
	"net/http"
)

type (
	RestServer struct {
		host    string
		port    int
		ssl     bool
		timeout int
	}
)

func NewRestServer(host string, port ...int) *RestServer {
	ret := &RestServer{}
	ret.host = host
	ret.port = 80
	ret.ssl = false
	ret.timeout = 0
	if len(port) > 0 {
		ret.port = port[0]
	}
	return ret
}

func NewSSLRestServer(host string, port ...int) *RestServer {
	ret := &RestServer{}
	ret.host = host
	ret.port = 443
	ret.ssl = true
	ret.timeout = 0
	if len(port) > 0 {
		ret.port = port[0]
	}
	return ret
}

func (rs *RestServer) JsonAPI(uri string, f ...func(*http.Request) *http.Request) func() error {
	return nil
}
