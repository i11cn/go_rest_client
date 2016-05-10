package rc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type (
	JsonBodyProcess struct {
	}
)

func (j *JsonBodyProcess) Marshal(body interface{}, req *http.Request) (err error) {
	if body != nil {
		var d []byte
		if s, ok := body.(string); ok {
			d = []byte(s)
		} else if sp, ok := body.(*string); ok {
			d = []byte(*sp)
		} else {
			d, err = json.Marshal(body)
			if err != nil {
				return err
			}
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(d))
		req.ContentLength = int64(len(d))
		req.Header.Set("Content-Type", "application/json;charset=utf-8")
	}
	return nil
}

func (j *JsonBodyProcess) Unmarshal(resp *http.Response, body []byte, obj interface{}) (err error) {
	return json.Unmarshal(body, obj)
}
