package rc

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type (
	Call interface {
		Invoke(*Request, ...interface{}) (*http.Response, error)
	}
	Caller func(*Request, ...interface{}) (*http.Response, error)

	BodyMarshal interface {
		Marshal(interface{}, *http.Request) error
	}
	BodyUnmarshal interface {
		Unmarshal(*http.Response, []byte, interface{}) error
	}

	RestClient struct {
		Host         string
		Port         int
		BaseUri      string
		SSL          bool
		body_marshal BodyMarshal
	}
	Request struct {
		Query  map[string]interface{}
		Body   interface{}
		Result interface{}
	}
	call_remote struct {
		client *RestClient
		method string
		uri    string
	}
)

var (
	g_cert_verify                      = false
	g_default_body_marshal BodyMarshal = &JsonBodyProcess{}
)

func (c *call_remote) Invoke(r *Request, args ...interface{}) (*http.Response, error) {
	uri := fmt.Sprintf(c.uri, args...)
	client := &http.Client{}
	if !g_cert_verify {
		client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}
	req, err := http.NewRequest(c.method, uri, nil)
	if err != nil {
		return nil, err
	}
	if len(r.Query) > 0 {
		buf := bytes.NewBufferString("")
		empty := true
		for k, v := range r.Query {
			if !empty {
				buf.WriteString("&")
			}
			buf.WriteString(k)
			buf.WriteString("=")
			buf.WriteString(fmt.Sprint(v))
		}
		req.URL.RawQuery = buf.String()
	}
	c.client.body_marshal.Marshal(r.Body, req)
	resp, err := client.Do(req)
	if err == nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil && r.Result != nil {
			ct := resp.Header.Get("Content-Type")
			switch {
			case strings.HasPrefix(strings.ToLower(ct), "application/json"):
				return resp, (&JsonBodyProcess{}).Unmarshal(resp, body, r.Result)
			case strings.HasPrefix(strings.ToLower(ct), "application/x-www-form-urlencoded"):
				return resp, (&FormBodyProcess{}).Unmarshal(resp, body, r.Result)
			}
		}
	}
	return resp, err
}

func VerifyCert(verify bool) {
	g_cert_verify = verify
}

func SetDefaultBodyMarshal(bm BodyMarshal) {
	g_default_body_marshal = bm
}

func NewClient(host string, port int, uri string) *RestClient {
	return &RestClient{Host: host, Port: port, BaseUri: uri, body_marshal: g_default_body_marshal}
}

func (rc *RestClient) SetBodyMarshal(bm BodyMarshal) {
	rc.body_marshal = bm
}

func (rc *RestClient) GetCaller(method, uri string) Caller {
	var buf bytes.Buffer
	if rc.SSL {
		buf.WriteString("https://")
		buf.WriteString(rc.Host)
		if rc.Port != 0 && rc.Port != 443 {
			buf.WriteByte(':')
			buf.WriteString(strconv.Itoa(rc.Port))
		}
	} else {
		buf.WriteString("http://")
		buf.WriteString(rc.Host)
		if rc.Port != 0 && rc.Port != 443 {
			buf.WriteByte(':')
			buf.WriteString(strconv.Itoa(rc.Port))
		}
	}
	if len(rc.BaseUri) > 0 && []byte(rc.BaseUri)[0] != '/' {
		buf.WriteByte('/')
	}
	buf.WriteString(rc.BaseUri)
	if len(uri) > 0 && []byte(uri)[0] != '/' {
		buf.WriteByte('/')
	}
	buf.WriteString(uri)
	return func(r *Request, args ...interface{}) (*http.Response, error) {
		cr := &call_remote{rc, method, buf.String()}
		return cr.Invoke(r, args...)
	}
}

func (rc *RestClient) Get(uri string, obj interface{}, args ...interface{}) (*http.Response, error) {
	return rc.GetCaller("GET", uri)(&Request{map[string]interface{}{}, nil, obj}, args...)
}

func (rc *RestClient) Post(uri string, body interface{}, obj interface{}, args ...interface{}) (*http.Response, error) {
	return rc.GetCaller("POST", uri)(&Request{map[string]interface{}{}, body, obj}, args...)
}

func (rc *RestClient) Put(uri string, body interface{}, obj interface{}, args ...interface{}) (*http.Response, error) {
	return rc.GetCaller("PUT", uri)(&Request{map[string]interface{}{}, body, obj}, args...)
}

func (rc *RestClient) Delete(uri string, body interface{}, obj interface{}, args ...interface{}) (*http.Response, error) {
	return rc.GetCaller("DELETE", uri)(&Request{map[string]interface{}{}, body, obj}, args...)
}

func (rc *RestClient) Option(uri string, body interface{}, obj interface{}, args ...interface{}) (*http.Response, error) {
	return rc.GetCaller("OPTION", uri)(&Request{map[string]interface{}{}, body, obj}, args...)
}

func (rc *RestClient) Head(uri string, args ...interface{}) (*http.Response, error) {
	return rc.GetCaller("HEAD", uri)(&Request{map[string]interface{}{}, nil, nil}, args...)
}

func (rc *RestClient) Trace(uri string, args ...interface{}) (*http.Response, error) {
	return rc.GetCaller("TRACE", uri)(&Request{map[string]interface{}{}, nil, nil}, args...)
}

func Get(uri string, obj interface{}, args ...interface{}) (*http.Response, error) {
	return rc.GetCaller("GET", uri)(&Request{map[string]interface{}{}, nil, obj}, args...)
}

func Post(uri string, body interface{}, obj interface{}, args ...interface{}) (*http.Response, error) {
	return rc.GetCaller("POST", uri)(&Request{map[string]interface{}{}, body, obj}, args...)
}

func Put(uri string, body interface{}, obj interface{}, args ...interface{}) (*http.Response, error) {
	return rc.GetCaller("PUT", uri)(&Request{map[string]interface{}{}, body, obj}, args...)
}

func Delete(uri string, body interface{}, obj interface{}, args ...interface{}) (*http.Response, error) {
	return rc.GetCaller("DELETE", uri)(&Request{map[string]interface{}{}, body, obj}, args...)
}

func Option(uri string, body interface{}, obj interface{}, args ...interface{}) (*http.Response, error) {
	return rc.GetCaller("OPTION", uri)(&Request{map[string]interface{}{}, body, obj}, args...)
}

func Head(uri string, args ...interface{}) (*http.Response, error) {
	return rc.GetCaller("HEAD", uri)(&Request{map[string]interface{}{}, nil, nil}, args...)
}

func Trace(uri string, args ...interface{}) (*http.Response, error) {
	return rc.GetCaller("TRACE", uri)(&Request{map[string]interface{}{}, nil, nil}, args...)
}
