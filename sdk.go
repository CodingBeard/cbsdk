package cbsdk

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	SdkMethodPost   = 1
	SdkMethodGet    = 2
	SdkMethodPut    = 3
	SdkMethodDelete = 4
)

type Sdk struct {
	timeout       time.Duration
	debug         bool
	httpMethod    uint8
	defaultHeader map[string]string
	header        map[string]string
	host          string
	uri           string
	get           []GetParam
	body          interface{}
	lock          *sync.Mutex
}

type GetParam struct {
	Key,
	Value string
}

type SdkRequest interface {
	GetDebug() bool
	GetHttpMethod() uint8
	GetHeader() map[string]string
	GetUri() string
	GetGet() []GetParam
	GetBody() interface{}
	GetTimeout() time.Duration
}

func New(host string, defaultHeaders map[string]string) *Sdk {
	return &Sdk{
		host:          host,
		defaultHeader: defaultHeaders,
		lock:          &sync.Mutex{},
	}
}

func (a *Sdk) AddDefaultHeader(key, value string) {
	a.defaultHeader[key] = value
}

func (a *Sdk) Create(request SdkRequest) (string, error) {
	a.lock.Lock()

	a.timeout = request.GetTimeout()
	a.debug = request.GetDebug()
	a.httpMethod = request.GetHttpMethod()
	a.header = request.GetHeader()
	a.uri = request.GetUri()
	a.get = request.GetGet()
	a.body = request.GetBody()

	var response string
	var e error

	switch a.httpMethod {
	case SdkMethodPost:
		response, e = a.makePostRequest()
	case SdkMethodGet:
		response, e = a.makeGetRequest()
	case SdkMethodPut:
		response, e = a.makePutRequest()
	case SdkMethodDelete:
		response, e = a.makeDeleteRequest()
	default:
		panic("specify a valid http method")
	}

	a.lock.Unlock()

	return response, e
}

func (a *Sdk) addGetParams(request *http.Request) {
	query := request.URL.Query()
	for _, param := range a.get {
		query.Add(param.Key, param.Value)
	}
	request.URL.RawQuery = query.Encode()
}

func (a *Sdk) addHeaders(request *http.Request) {
	for key, value := range a.header {
		request.Header.Set(key, value)
	}
}

func (a *Sdk) makePostRequest() (string, error) {
	url := strings.Join([]string{a.host, a.uri}, "")

	client := &http.Client{
		Timeout: a.timeout,
	}

	var requestBody []byte
	var e error
	_, isStringMap := a.body.(map[string]interface{})
	_, isString := a.body.(string)
	if isStringMap {
		requestBody, e = json.Marshal(a.body)
		if e != nil {
			return "", e
		}
	} else if isString {
		requestBody = []byte(a.body.(string))
	}

	request, e := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if e != nil {
		return "", e
	}

	a.addGetParams(request)
	a.addHeaders(request)

	request.Header.Set("Content-Type", "application/json")

	response, e := client.Do(request)
	if e != nil {
		return "", e
	}

	defer response.Body.Close()
	body, e := ioutil.ReadAll(response.Body)
	if e != nil {
		return "", e
	}

	return string(body), nil
}

func (a *Sdk) makeGetRequest() (string, error) {
	url := strings.Join([]string{a.host, a.uri}, "")

	client := &http.Client{
		Timeout: a.timeout,
	}

	request, e := http.NewRequest("GET", url, nil)
	if e != nil {
		return "", e
	}

	a.addGetParams(request)
	a.addHeaders(request)

	response, e := client.Do(request)
	if e != nil {
		return "", e
	}

	defer response.Body.Close()
	body, e := ioutil.ReadAll(response.Body)
	if e != nil {
		return "", e
	}

	return string(body), nil
}

func (a *Sdk) makePutRequest() (string, error) {
	url := strings.Join([]string{a.host, a.uri}, "")

	client := &http.Client{
		Timeout: a.timeout,
	}

	var requestBody []byte
	var e error
	_, isStringMap := a.body.(map[string]interface{})
	_, isString := a.body.(string)
	if isStringMap {
		requestBody, e = json.Marshal(a.body)
		if e != nil {
			return "", e
		}
	} else if isString {
		requestBody = []byte(a.body.(string))
	}

	request, e := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
	if e != nil {
		return "", e
	}

	a.addGetParams(request)
	a.addHeaders(request)

	request.Header.Set("Content-Type", "application/json")

	response, e := client.Do(request)
	if e != nil {
		return "", e
	}

	defer response.Body.Close()
	body, e := ioutil.ReadAll(response.Body)
	if e != nil {
		return "", e
	}

	return string(body), nil
}

func (a *Sdk) makeDeleteRequest() (string, error) {
	url := strings.Join([]string{a.host, a.uri}, "")

	client := &http.Client{
		Timeout: a.timeout,
	}

	request, e := http.NewRequest("DELETE", url, nil)
	if e != nil {
		return "", e
	}

	a.addGetParams(request)
	a.addHeaders(request)

	response, e := client.Do(request)
	if e != nil {
		return "", e
	}

	defer response.Body.Close()
	body, e := ioutil.ReadAll(response.Body)
	if e != nil {
		return "", e
	}

	return string(body), nil
}
