package grc

import (
	"fmt"
	"net/http"
)

type (
	RestServer struct {
		host string
		port int
		ssl  bool
		url  string
	}

	API interface {
		Run(...interface{}) (Response, error)
	}
)

func NewRestServer(host string, port ...int) *RestServer {
	ret := &RestServer{}
	ret.host = host
	ret.port = 80
	ret.ssl = false
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
	if len(port) > 0 {
		ret.port = port[0]
	}
	return ret
}

func (rs *RestServer) get_url() string {
	if rs.url == "" {
		if rs.ssl {
			if rs.port == 443 || rs.port == 0 {
				return fmt.Sprintf("https://%s", rs.host)
			} else {
				return fmt.Sprintf("https://%s:%d", rs.host, rs.port)
			}

		} else {
			if rs.port == 80 || rs.port == 0 {
				return fmt.Sprintf("http://%s", rs.host)
			} else {
				return fmt.Sprintf("http://%s:%d", rs.host, rs.port)
			}
		}
	}
	return rs.url
}

func (rs *RestServer) JsonAPI(method, uri string, body bool, f ...func(*http.Request)) API {
	return new_json_api(rs, method, uri, body, f...)
}
