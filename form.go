package rc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
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
			for k, v := range b {
				buf.WriteString(fmt.Sprintf("%s=%v", k, v))
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
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}
	return nil
}

func (j *FormBodyProcess) Unmarshal(body []byte, obj interface{}) (err error) {
	obj = string(body)
	return nil
}
