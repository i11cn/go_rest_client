package grc

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type (
	Body interface {
		Length() int64
		Stream() io.Reader
		Bytes() ([]byte, error)
		String() (string, error)
		Json(interface{}) error
		Xml(interface{}) error
	}

	Response interface {
		StatusCode() int
		Status() string
		Header() http.Header
		Body() Body
	}

	grc_body struct {
		body   io.Reader
		data   []byte
		length int64
	}

	grc_response struct {
		req  *http.Request
		resp *http.Response
		body Body
	}
)

func new_grc_body(resp *http.Response, body io.Reader) Body {
	ret := &grc_body{}
	ret.body = body
	ret.length = resp.ContentLength
	return ret
}

func (gb *grc_body) Length() int64 {
	return gb.length
}

func (gb *grc_body) Stream() io.Reader {
	return gb.body
}

func (gb *grc_body) Bytes() (ret []byte, err error) {
	if gb.data == nil {
		ret, err = ioutil.ReadAll(gb.body)
		if err == nil {
			gb.data = ret
		}
		return
	}
	return gb.data, nil
}

func (gb *grc_body) String() (string, error) {
	d, err := gb.Bytes()
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func (gb *grc_body) Json(obj interface{}) error {
	d, err := gb.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(d, obj)
}

func (gb *grc_body) Xml(obj interface{}) error {
	d, err := gb.Bytes()
	if err != nil {
		return err
	}
	return xml.Unmarshal(d, obj)
}

func new_grc_response(resp *http.Response, ungzip bool) Response {
	ret := &grc_response{}
	ret.req = resp.Request
	ret.resp = resp
	var rd io.Reader = resp.Body
	if ungzip && len(resp.Header.Get("Content-Encoding")) > 0 {
		switch strings.ToUpper(resp.Header.Get("Content-Encoding")) {
		case "GZIP":
			tmp, err := gzip.NewReader(resp.Body)
			if err == nil {
				rd = tmp
			}

		case "DEFLATE":
			rd = flate.NewReader(resp.Body)
		}
	}
	ret.body = new_grc_body(resp, rd)
	return ret
}

func (gr *grc_response) StatusCode() int {
	return gr.resp.StatusCode
}

func (gr *grc_response) Status() string {
	return gr.resp.Status
}

func (gr *grc_response) Header() http.Header {
	return gr.resp.Header
}

func (gr *grc_response) Body() Body {
	return gr.body
}
