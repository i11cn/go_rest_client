package rc

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type (
	XmlBodyProcess struct {
	}
)

func (j *XmlBodyProcess) Marshal(body interface{}, req *http.Request) (err error) {
	if body != nil {
		var d []byte
		if s, ok := body.(string); ok {
			d = []byte(s)
		} else if sp, ok := body.(*string); ok {
			d = []byte(*sp)
		} else {
			d, err = xml.Marshal(body)
			if err != nil {
				return err
			}
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(d))
		req.ContentLength = int64(len(d))
		req.Header.Set("Content-Type", "application/xml;charset=utf-8")
	}
	return nil
}

func (j *XmlBodyProcess) Unmarshal(resp *http.Response, body []byte, obj interface{}) (err error) {
	return xml.Unmarshal(body, obj)
}
