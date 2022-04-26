package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	common "bitbucket.org/MarkEdwardTresidder/micro-common"
	"github.com/gin-gonic/gin"
)

type IHttp interface {
	Invoke(log *common.MicroLog, methodType, url string, jsonForm []byte, vs ...interface{}) ([]byte, error)
	SetHeaders(http *http.Header, c *gin.Context)
	GenerateJson(vs interface{}) ([]byte, error)
}

type Http struct {
}

func NewHttp() IHttp {
	return &Http{}
}

type QueryParam map[string]interface{}
type Param map[string]interface{}

type param struct {
	url.Values
}

func (p *param) getValues() url.Values {
	if p.Values == nil {
		p.Values = make(url.Values)
	}
	return p.Values
}

func (p *param) Adds(m map[string]interface{}) {
	if len(m) == 0 {
		return
	}
	vs := p.getValues()
	for k, v := range m {
		vs.Add(k, fmt.Sprint(v))
	}
}
func (p *param) Empty() bool {
	return p.Values == nil
}

func (h *Http) Invoke(log *common.MicroLog, methodType, rawUrl string, jsonForm []byte, vs ...interface{}) ([]byte, error) {

	client := &http.Client{}
	var req *http.Request
	var err error

	if rawUrl == "" {
		log.Message(`url not specified`)
		return nil, errors.New(`url not specified`)
	}
	switch m := methodType; m {
	case http.MethodPost:
		req, err = http.NewRequest(http.MethodPost, rawUrl, bytes.NewBuffer(jsonForm))
	case http.MethodPut:
		req, err = http.NewRequest(http.MethodPut, rawUrl, bytes.NewBuffer(jsonForm))
	case http.MethodPatch:
		req, err = http.NewRequest(http.MethodPatch, rawUrl, bytes.NewBuffer(jsonForm))
	case http.MethodGet:
		req, err = http.NewRequest(http.MethodGet, rawUrl, nil)
	default:
		log.Message(`Invalid Method Type`)
		return nil, errors.New(`Invalid Method Type`)
	}

	if err != nil {
		log.Message("Error while building Request" + err.Error())
		return nil, err
	}

	var queryParam param

	for _, v := range vs {
		switch vv := v.(type) {
		case http.Header:
			for key, values := range vv {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}
		case QueryParam:
			queryParam.Adds(vv)
		case error:
			return nil, vv
		default:
			log.Message(`Invalid Argument`)
			return nil, errors.New(`Invalid Argument`)
		}
	}

	if !queryParam.Empty() {
		paramStr := queryParam.Encode()
		if strings.IndexByte(rawUrl, '?') == -1 {
			rawUrl = rawUrl + "?" + paramStr
		} else {
			rawUrl = rawUrl + "&" + paramStr
		}
	}

	u, err := url.Parse(rawUrl)
	if err != nil {
		log.Message(`Error while parsing URL` + err.Error())
		return nil, err
	}
	req.URL = u

	resp, err := client.Do(req)

	if err != nil {
		log.Message("Error while invoking service" + err.Error())
		return nil, err
	}

	apiData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Message("Error while reading Response" + err.Error())
		return nil, err
	}

	resp.Body.Close()
	return apiData, nil
}

func (h *Http) SetHeaders(header *http.Header, c *gin.Context) {
	header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	header.Set("X-Tenant-Id", c.Request.Header.Get("X-Tenant-Id"))
	header.Set("X-User-Id", c.Request.Header.Get("X-User-Id"))
	if c.Request.Header.Get("X-Auth-Type") == "vendor" {
		header.Set("X-Reference-Id", c.Request.Header.Get("X-Reference-Id"))
	}
	header.Set("X-B3-Traceid", c.Request.Header.Get("X-B3-Traceid"))
	header.Set("X-B3-Spanid", c.Request.Header.Get("X-B3-Spanid"))
	header.Set("X-Name", c.Request.Header.Get("X-Name"))
	header.Set("Content-Type", "application/json")
}

func (h *Http) GenerateJson(vs interface{}) ([]byte, error) {
	json, err := json.Marshal(vs)
	if err != nil {
		return nil, err
	}
	return json, nil
}
