package grc

import (
	"errors"
)

type (
	APIPool interface {
		Register(string, API)
		Call(string, ...interface{}) (Response, error)
	}

	grc_api_pool struct {
		pool map[string]API
	}
)

func (ap *grc_api_pool) Register(name string, api API) {
	ap.pool[name] = api
}

func (ap *grc_api_pool) Call(name string, obj ...interface{}) (Response, error) {
	if api, exist := ap.pool[name]; exist {
		return api.Run(obj...)
	} else {
		return nil, errors.New("API " + name + " not exist")
	}
}

func NewAPIPool() APIPool {
	ret := &grc_api_pool{}
	ret.pool = make(map[string]API)
	return ret
}
