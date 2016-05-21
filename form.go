package rc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type (
	FormBodyProcess struct {
	}
)

func (j *FormBodyProcess) Marshal(body interface{}, req *http.Request) (err error) {
	if body != nil {
		buf := new(bytes.Buffer)
		valid := true
		switch b := body.(type) {
		case map[string]interface{}:
			count := 0
			for k, v := range b {
				if count > 0 {
					buf.WriteString("&")
				}
				count++
				value := fmt.Sprintf("%v", v)
				buf.WriteString(fmt.Sprintf("%s=%s", k, url.QueryEscape(value)))
			}
		case map[string]string:
			count := 0
			for k, v := range b {
				if count > 0 {
					buf.WriteString("&")
				}
				count++
				buf.WriteString(fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
			}
		case map[string][]string:
			count := 0
			for k, vs := range b {
				for _, v := range vs {
					if count > 0 {
						buf.WriteString("&")
					}
					count++
					buf.WriteString(fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
				}
			}
		case string:
			buf.WriteString(b)
		case *string:
			buf.WriteString(*b)
		default:
			valid = false
		}
		if valid {
			req.Body = ioutil.NopCloser(buf)
			req.ContentLength = int64(buf.Len())
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}
	return nil
}

func (j *FormBodyProcess) Unmarshal(resp *http.Response, body []byte, obj interface{}) (err error) {
	switch o := obj.(type) {
	case map[string][]string:
		o = resp.Request.PostForm
	case *string:
		*o = string(body)
	}
	return nil
}
