package utils

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	qbc "github.com/rskvp/qb-core"
	"github.com/valyala/fasthttp"
)

type ResponseData struct {
	StatusCode int
	Body       []byte
	Header     map[string]string
}

func NewResponseDataEmpty() *ResponseData {
	instance := new(ResponseData)
	return instance
}

func NewResponseData(res *fasthttp.Response) *ResponseData {
	instance := new(ResponseData)
	instance.StatusCode = res.StatusCode()
	instance.Body = res.Body()
	instance.Header = responseHeaderToMap(res.Header.String())

	return instance
}

func (instance *ResponseData) String() string {
	m := map[string]interface{}{
		"status": instance.StatusCode,
		"header": instance.Header,
		"body":   string(instance.Body),
	}
	return qbc.JSON.Stringify(m)
}

func (instance *ResponseData) BodyAsMap() map[string]interface{} {
	var m map[string]interface{}
	err := qbc.JSON.Read(instance.Body, &m)
	if nil == err {
		return m
	}
	return map[string]interface{}{}
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func HttpHeaderToMap(header http.Header) map[string]string {
	response := make(map[string]string)
	for k, v := range header {
		response[k] = strings.Join(v, ",")
	}
	return response
}

func Params(ctx *fiber.Ctx) map[string]interface{} {
	// get all body params
	var response map[string]interface{}
	_ = qbc.JSON.Read(ctx.Body(), &response)
	if nil == response {
		response = map[string]interface{}{}
	}

	// try add form params
	if form, err := ctx.MultipartForm(); nil == err && nil != form && nil != form.Value {
		for k, v := range form.Value {
			response[k] = v
		}
	}

	// url query
	if path := ctx.OriginalURL(); len(path) > 0 {
		uri, err := url.Parse(path)
		if nil == err {
			query := uri.Query()
			if nil != query && len(query) > 0 {
				for k, v := range query {
					if len(v) == 1 {
						response[k] = v[0]
					} else {
						response[k] = v
					}
				}
			}
		}
	}

	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func responseHeaderToMap(text string) map[string]string {
	response := map[string]string{}
	if len(text) > 0 {
		tokens := qbc.Strings.Split(text, "\n\r")
		for _, v := range tokens {
			pair := qbc.Strings.Split(v, ":")
			if len(pair) > 1 {
				response[pair[0]] = strings.Join(pair[1:], ":")
			} else {
				response["http"] = v
			}
		}
	}
	return response
}
