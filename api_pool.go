package grc

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type (
	APICaller func(...interface{}) (Response, error)

	tag_info struct {
		host   string
		port   int
		ssl    bool
		method string
		uri    string
		body   bool
	}
)

var (
	grc_server_map map[string]*RestServer
)

func init() {
	grc_server_map = make(map[string]*RestServer)
}

func (ti *tag_info) get_rest_server() *RestServer {
	var url string
	if ti.ssl {
		if ti.port == 0 || ti.port == 443 {
			url = fmt.Sprintf("https://%s", ti.host)
		} else {
			url = fmt.Sprintf("https://%s:%d", ti.host, ti.port)
		}
	} else {
		if ti.port == 0 || ti.port == 80 {
			url = fmt.Sprintf("http://%s", ti.host)
		} else {
			url = fmt.Sprintf("http://%s:%d", ti.host, ti.port)
		}
	}
	if s, exist := grc_server_map[url]; exist {
		return s
	}
	var ret *RestServer
	if ti.ssl {
		ret = NewSSLRestServer(ti.host, ti.port)
	} else {
		ret = NewRestServer(ti.host, ti.port)
	}
	grc_server_map[url] = ret
	return ret
}

func parse_struct_tag(tag reflect.StructTag) (*tag_info, error) {
	info := &tag_info{}
	info.host = tag.Get("grc_host")
	if info.host == "" {
		return nil, errors.New("Tag中必须设置grc_host")
	}
	info.uri = tag.Get("grc_uri")
	if info.uri == "" {
		return nil, errors.New("Tag中必须设置grc_uri")
	}
	if tag.Get("grc_flags") == "" {
		return nil, errors.New("至少要在grc_flags中指定调用的Method")
	}
	flags := strings.Split(tag.Get("grc_flags"), ",")
	for _, flag := range flags {
		flag = strings.ToUpper(flag)
		switch flag {
		case "SSL":
			info.ssl = true

		case "BODY":
			info.body = true

		case "GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "TRACE":
			info.method = flag

		default:
			if num, err := regexp.MatchString("^\\d+$", flag); err == nil && num {
				info.port, _ = strconv.Atoi(flag)
			}
		}
	}
	if info.method == "" {
		return nil, errors.New("至少要在grc_flags中指定调用的Method")
	}
	return info, nil
}

// Host , port , ssl, method, uri, body

func process_struct_field(v reflect.Value, t reflect.StructField) error {
	info, err := parse_struct_tag(t.Tag)
	if err != nil {
		return err
	}
	server := info.get_rest_server()
	api := server.JsonAPI(info.method, info.uri, info.body)
	wrapper := func(obj ...interface{}) (Response, error) {
		return api.Run(obj...)
	}
	v.Set(reflect.ValueOf(wrapper))
	return nil
}

func ParsePool(obj interface{}) error {
	v := reflect.ValueOf(obj)
	t := v.Type()
	if t.Kind() != reflect.Ptr {
		return errors.New("非指针类型的对象设置无效")
	}
	v = v.Elem()
	t = v.Type()
	if t.Kind() != reflect.Struct {
		return errors.New("非结构指针类型的对象设置无效")
	}
	for i := 0; i < t.NumField(); i++ {
		f := v.Field(i)
		if f.Type().String() == "grc.APICaller" {
			if err := process_struct_field(f, t.Field(i)); err != nil {
				return err
			}
		}
	}
	return nil
}
