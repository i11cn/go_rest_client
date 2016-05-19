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
		Invoke(*Request, ...interface{}) error
	}
	Caller func(*Request, ...interface{}) error

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

func (c *call_remote) Invoke(r *Request, args ...interface{}) error {
	uri := fmt.Sprintf(c.uri, args...)
	client := &http.Client{}
	if !g_cert_verify {
		client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}
	req, err := http.NewRequest(c.method, uri, nil)
	if err != nil {
		return err
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
			case strings.HasPrefix(ct, "application/json"):
				return (&JsonBodyProcess{}).Unmarshal(resp, body, r.Result)
			case strings.HasPrefix(ct, "application/x-www-form-urlencoded"):
				return (&FormBodyProcess{}).Unmarshal(resp, body, r.Result)
			}
		}
	}
	return err

	return nil
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
	return func(r *Request, args ...interface{}) error {
		cr := &call_remote{rc, method, buf.String()}
		return cr.Invoke(r, args...)
	}
}

/*
func Get(host string, port int, uri string, obj interface{}) error {
	c := &RestClient{Method: "GET", Host: host, Port: port, Uri: uri, body_marshal: g_default_body_marshal}
	return c.Do(nil, obj)
}

func Post(host string, port int, uri string, body, obj interface{}) error {
	c := &RestClient{Method: "POST", Host: host, Port: port, Uri: uri, body_marshal: g_default_body_marshal}
	r := &Request{map[string]interface{}{}, body}
	return c.Do(r, obj)
}

func Put(host string, port int, uri string, body, obj interface{}) error {
	c := &RestClient{Method: "PUT", Host: host, Port: port, Uri: uri, body_marshal: g_default_body_marshal}
	r := &Request{map[string]interface{}{}, body}
	return c.Do(r, obj)
}

func Delete(host string, port int, uri string, body, obj interface{}) error {
	c := &RestClient{Method: "DELETE", Host: host, Port: port, Uri: uri, body_marshal: g_default_body_marshal}
	r := &Request{map[string]interface{}{}, body}
	return c.Do(r, obj)
}

func Option(host string, port int, uri string, body, obj interface{}) error {
	c := &RestClient{Method: "OPTION", Host: host, Port: port, Uri: uri, body_marshal: g_default_body_marshal}
	r := &Request{map[string]interface{}{}, body}
	return c.Do(r, obj)
}

func Head(host string, port int, uri string, body, obj interface{}) error {
	c := &RestClient{Method: "HEAD", Host: host, Port: port, Uri: uri, body_marshal: g_default_body_marshal}
	r := &Request{map[string]interface{}{}, body}
	return c.Do(r, obj)
}

func Patch(host string, port int, uri string, body, obj interface{}) error {
	c := &RestClient{Method: "PATCH", Host: host, Port: port, Uri: uri, body_marshal: g_default_body_marshal}
	r := &Request{map[string]interface{}{}, body}
	return c.Do(r, obj)
}

func Trace(host string, port int, uri string, body, obj interface{}) error {
	c := &RestClient{Method: "TRACE", Host: host, Port: port, Uri: uri, body_marshal: g_default_body_marshal}
	r := &Request{map[string]interface{}{}, body}
	return c.Do(r, obj)
}*/
