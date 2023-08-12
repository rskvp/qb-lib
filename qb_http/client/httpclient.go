package client

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	qbc "github.com/rskvp/qb-core"
	httputils "github.com/rskvp/qb-lib/qb_http/utils"
	"github.com/valyala/fasthttp"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const methodGet = "GET"
const methodPost = "POST"
const methodPut = "PUT"
const methodDelete = "DELETE"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type HttpClient struct {
	client *fasthttp.Client
	header map[string]string
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewHttpClient() *HttpClient {
	instance := new(HttpClient)
	instance.header = make(map[string]string)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpClient) AddHeader(key, value string) {
	if nil != instance && nil != instance.header {
		instance.header[key] = value
	}
}

func (instance *HttpClient) RemoveHeader(key string) {
	if nil != instance && nil != instance.header {
		delete(instance.header, key)
	}
}

func (instance *HttpClient) Get(url string) (*httputils.ResponseData, error) {
	// return instance.GetTimeout(url, time.Second*15)
	return instance.do(methodGet, url, nil, time.Second*15)
}

func (instance *HttpClient) GetTimeout(url string, timeout time.Duration) (*httputils.ResponseData, error) {
	// return instance.get(url, timeout)
	return instance.do(methodGet, url, nil, timeout)
}

func (instance *HttpClient) Post(url string, body interface{}) (*httputils.ResponseData, error) {
	// return instance.GetTimeout(url, time.Second*15)
	return instance.do(methodPost, url, body, time.Second*15)
}

func (instance *HttpClient) PostTimeout(url string, body interface{}, timeout time.Duration) (*httputils.ResponseData, error) {
	// return instance.get(url, timeout)
	return instance.do(methodPost, url, body, timeout)
}

func (instance *HttpClient) Put(url string, body interface{}) (*httputils.ResponseData, error) {
	// return instance.GetTimeout(url, time.Second*15)
	return instance.do(methodPut, url, body, time.Second*15)
}

func (instance *HttpClient) PutTimeout(url string, body interface{}, timeout time.Duration) (*httputils.ResponseData, error) {
	// return instance.get(url, timeout)
	return instance.do(methodPut, url, body, timeout)
}

func (instance *HttpClient) Delete(url string, body interface{}) (*httputils.ResponseData, error) {
	// return instance.GetTimeout(url, time.Second*15)
	return instance.do(methodDelete, url, body, time.Second*15)
}

func (instance *HttpClient) DeleteTimeout(url string, body interface{}, timeout time.Duration) (*httputils.ResponseData, error) {
	// return instance.get(url, timeout)
	return instance.do(methodDelete, url, body, timeout)
}

func (instance *HttpClient) Upload(url string, filename, optParamName string, optParams map[string]interface{}) (*httputils.ResponseData, error) {
	return instance.UploadTimeout(url, filename, optParamName, optParams, time.Second*120)
}

func (instance *HttpClient) UploadTimeout(url string, filename, optParamName string, optParams map[string]interface{}, timeout time.Duration) (*httputils.ResponseData, error) {
	if nil == optParams {
		optParams = make(map[string]interface{})
	}
	if len(optParamName) == 0 {
		optParamName = "file"
	}

	// open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	part, err := writer.CreateFormFile(optParamName, filepath.Base(filename))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	for key, val := range optParams {
		_ = writer.WriteField(key, qbc.Convert.ToString(val))
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if nil != instance.header {
		// fmt.Println(req.Header.String())
		for k, v := range instance.header {
			req.Header.Set(k, v)
		}
	}

	client := &http.Client{}
	client.Timeout = timeout
	resp, err := client.Do(req)
	if nil != err {
		return nil, err
	}
	defer resp.Body.Close()

	response := httputils.NewResponseDataEmpty()
	response.Body, _ = ioutil.ReadAll(resp.Body)
	response.StatusCode = resp.StatusCode
	response.Header = httputils.HttpHeaderToMap(resp.Header)
	return response, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpClient) init() *fasthttp.Client {
	if nil == instance.client {
		instance.client = new(fasthttp.Client)
		instance.client.Name = "qb_http_client/1.0"
	}
	return instance.client
}

func (instance *HttpClient) do(method string, uri string, reqBody interface{}, timeout time.Duration) (*httputils.ResponseData, error) {
	var err error

	// request
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(method)
	req.SetRequestURI(uri)

	if nil != instance.header {
		// fmt.Println(req.Header.String())
		for k, v := range instance.header {
			req.Header.Set(k, v)
		}
	}

	if nil != reqBody {
		if v, b := reqBody.(string); b {
			req.SetBodyString(v)
		} else if v, b := reqBody.([]byte); b {
			req.SetBody(v)
		} else if v, b := reqBody.([]uint8); b {
			req.SetBody(v)
		} else if v, b := reqBody.(map[string]interface{}); b {
			req.SetBodyString(qbc.JSON.Stringify(v))
		}
	}

	// response
	res := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	client := instance.init()
	if isHttps(uri) {
		// tlsConfig := &tls.Config{}
		// client.TLSConfig = tlsConfig
	}
	err = client.DoTimeout(req, res, timeout)

	return httputils.NewResponseData(res), err
}

func (instance *HttpClient) get(uri string, timeout time.Duration) (statusCode int, body []byte, err error) {
	client := instance.init()
	return client.GetTimeout(nil, uri, timeout)
}

func isHttps(uri string) bool {
	return strings.Index(strings.ToLower(uri), "https") > -1
}
