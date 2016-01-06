package rc

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	RestClient struct {
		Method string
		Host   string
		Port   int
		Uri    string
		Query  map[string]interface{}
		Body   interface{}
		SSL    bool
	}
)

var g_cert_verify = false

func VerifyCert(verify bool) {
	g_cert_verify = verify
}

func NewClient(host string, port int, uri string, query map[string]interface{}, body interface{}) *RestClient {
	return &RestClient{Method: "GET", Host: host, Port: port, Uri: uri, Query: query, Body: body}
}

func (rc *RestClient) Do(obj interface{}) error {
	client := &http.Client{}
	if !g_cert_verify {
		client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}
	url := ""
	root := ""
	if len(rc.Uri) > 0 && []byte(rc.Uri)[0] != '/' {
		root = "/"
	}
	if rc.SSL {
		if rc.Port == 0 || rc.Port == 443 {
			url = fmt.Sprintf("https://%s%s%s", rc.Host, root, rc.Uri)
		} else {
			url = fmt.Sprintf("https://%s:%d%s%s", rc.Host, rc.Port, root, rc.Uri)
		}
	} else {
		if rc.Port == 0 || rc.Port == 80 {
			url = fmt.Sprintf("http://%s%s%s", rc.Host, root, rc.Uri)
		} else {
			url = fmt.Sprintf("http://%s:%d%s%s", rc.Host, rc.Port, root, rc.Uri)
		}
	}
	req, err := http.NewRequest(rc.Method, url, nil)
	if err != nil {
		return err
	}
	if len(rc.Query) > 0 {
		buf := bytes.NewBufferString("")
		empty := true
		for k, v := range rc.Query {
			if !empty {
				buf.WriteString("&")
			}
			buf.WriteString(k)
			buf.WriteString("=")
			buf.WriteString(fmt.Sprint(v))
		}
		req.URL.RawQuery = buf.String()
	}
	if rc.Body != nil {
		var d []byte
		if s, ok := obj.(string); ok {
			d = []byte(s)
		} else if sp, ok := obj.(*string); ok {
			d = []byte(*sp)
		} else {
			d, err = json.Marshal(obj)
			if err != nil {
				return err
			}
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(d))
		req.Header.Set("Content-Type", "application/json;charset=utf-8")
	}
	resp, err := client.Do(req)
	if err == nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil && obj != nil {
			return json.Unmarshal(body, obj)
		}
	}
	return err
}

func (rc *RestClient) Get(obj interface{}) error {
	rc.Method = "GET"
	return rc.Do(obj)
}

func (rc *RestClient) Post(obj interface{}) error {
	rc.Method = "POST"
	return rc.Do(obj)
}

func (rc *RestClient) Put(obj interface{}) error {
	rc.Method = "PUT"
	return rc.Do(obj)
}

func (rc *RestClient) Delete(obj interface{}) error {
	rc.Method = "DELETE"
	return rc.Do(obj)
}

func (rc *RestClient) Option(obj interface{}) error {
	rc.Method = "OPTION"
	return rc.Do(obj)
}

func (rc *RestClient) Head(obj interface{}) error {
	rc.Method = "HEAD"
	return rc.Do(obj)
}

func (rc *RestClient) Patch(obj interface{}) error {
	rc.Method = "PATCH"
	return rc.Do(obj)
}

func (rc *RestClient) Trace(obj interface{}) error {
	rc.Method = "TRACE"
	return rc.Do(obj)
}

func Get(host string, port int, uri string, obj interface{}) error {
	c := &RestClient{Method: "GET", Host: host, Port: port, Uri: uri}
	return c.Do(obj)
}

func Post(host string, port int, uri string, body, obj interface{}) error {
	c := &RestClient{Method: "POST", Host: host, Port: port, Uri: uri, Body: body}
	return c.Do(obj)
}

func Put(host string, port int, uri string, body, obj interface{}) error {
	c := &RestClient{Method: "PUT", Host: host, Port: port, Uri: uri, Body: body}
	return c.Do(obj)
}

func Delete(host string, port int, uri string, body, obj interface{}) error {
	c := &RestClient{Method: "DELETE", Host: host, Port: port, Uri: uri, Body: body}
	return c.Do(obj)
}

func Option(host string, port int, uri string, body, obj interface{}) error {
	c := &RestClient{Method: "OPTION", Host: host, Port: port, Uri: uri, Body: body}
	return c.Do(obj)
}

func Head(host string, port int, uri string, body, obj interface{}) error {
	c := &RestClient{Method: "HEAD", Host: host, Port: port, Uri: uri, Body: body}
	return c.Do(obj)
}

func Patch(host string, port int, uri string, body, obj interface{}) error {
	c := &RestClient{Method: "PATCH", Host: host, Port: port, Uri: uri, Body: body}
	return c.Do(obj)
}

func Trace(host string, port int, uri string, body, obj interface{}) error {
	c := &RestClient{Method: "TRACE", Host: host, Port: port, Uri: uri, Body: body}
	return c.Do(obj)
}
