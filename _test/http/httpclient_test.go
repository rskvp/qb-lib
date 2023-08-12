package _test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"

	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_http/httpclient"
	"github.com/valyala/fasthttp"
)

func TestSimple(t *testing.T) {

	client := new(httpclient.HttpClient)

	// Fetch google page via local proxy.
	fmt.Println("https://gianangelogeminiani.me")
	resp, err := client.Get("https://gianangelogeminiani.me")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	if resp.StatusCode != fasthttp.StatusOK {
		log.Fatalf("Unexpected status code: %d. Expecting %d", resp.StatusCode, fasthttp.StatusOK)
	}
	useResponseBody(resp.Body)

	// Fetch foobar page via local proxy. Reuse body buffer.
	fmt.Println("https://botika.ai/")
	resp, err = client.Get("https://botika.ai/")
	if resp.StatusCode != fasthttp.StatusOK {
		log.Fatalf("Unexpected status code: %d. Expecting %d", resp.StatusCode, fasthttp.StatusOK)
	}
	if err != nil {
		log.Fatalf("Error when loading google page through local proxy: %s", err)
	}
	useResponseBody(resp.Body)

}

func TestTinyURL(t *testing.T) {
	urlFull := "http://localhost:63343/ritiro_io_client/index.html?_ijt=qbk2r2ocijg43343og9ivnvr4o#!/02_viewer/menu/ee63b7f4-1766-487e-8762-3a2710320158/04eff121-d533-edc9-7fc2-ebc393895250"
	client := new(httpclient.HttpClient)
	callUrl := "http://tinyurl.com/api-create.php?url=" + url.QueryEscape(urlFull)
	response, err := client.Get(callUrl)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(string(response.Body))
}

func TestDownload(t *testing.T) {
	client := new(httpclient.HttpClient)

	file := "https://gianangelogeminiani.me/download/architecture.png"
	fmt.Println(file)
	resp, err := client.Get(file)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	if resp.StatusCode != fasthttp.StatusOK {
		log.Fatalf("Unexpected status code: %d. Expecting %d", resp.StatusCode, fasthttp.StatusOK)
	}

	fileName := "./downloads/architecture.png"
	_ = qbc.Paths.Mkdir(fileName)
	_, err = qbc.IO.WriteBytesToFile(resp.Body, fileName)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(fileName)
}

func TestUpload(t *testing.T) {
	uri := "https://httpbin.org/post" // "http://localhost:9090/api/v1/files/upload" //
	body := map[string]interface{}{
		"param1": "hello",
	}

	authToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiODFmMTYzNmUzOGZlYzVjZWZkMjlhNzQxMGU5MDcwMDciLCJwYXlsb2FkIjp7ImNvbmZpcm1lZCI6dHJ1ZSwidXNlcl9uYW1lIjoiZW5jLTF0Sm5SVHgwY3crdU82NG5ta2w3MmVHZ2kwVzlSWXF0Z3RNb3pUYVVLM1ovc0pjNjJEcjh5N2dNeUJ3S0pldVVXaWY4QWhIZyIsInVzZXJfcHN3IjoiZW5jLVJubTNhanlJeWNZdE9uZnIvTFNkWm1ZL0V6RTlmV3gwUHhDRGIwKzBhVksxaUE9PSIsInVzZXJfcHN3X3RpbWVzdGFtcCI6MTY1NTQ1ODA3MH0sInNlY3JldF90eXBlIjoiYWNjZXNzIiwiZXhwIjoxNjU1NTQxMDU0LCJqdGkiOiI1ZjJiYjRjNGIyNGZhMTk3ODcyNDgyZDI3ZGRhMTU0NSJ9.ZK5Takwzlpi6uFw8tp3kDhRiNyJMzDdJjO4VnOQzQsA"

	client := httpclient.NewHttpClient()
	client.AddHeader("Authorization", "Bearer "+authToken)
	response, err := client.Upload(uri, "file1.txt", "file", body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("RESPONSE: ", string(response.Body))
}

func TestPost(t *testing.T) {
	endpoint := "https://httpbin.org/post"
	body := map[string]interface{}{
		"param1": "hello",
	}
	client := httpclient.NewHttpClient()
	client.AddHeader("Content-Type", "application/json")
	response, err := client.Post(endpoint, body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(string(response.Body))
}

func TestPostFormData(t *testing.T) {
	endpoint := "https://httpbin.org/post"
	body := fmt.Sprintf("PARAM1=%s", "1234")
	client := httpclient.NewHttpClient()
	client.AddHeader("Accept-Encoding", "gzip")
	client.AddHeader("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Post(endpoint, body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(string(response.Body))
}

func TestFormDataNative(t *testing.T) {

	url := "https://httpbin.org/post"
	method := "POST"

	payload := strings.NewReader("PARAM1=1234")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}

func TestSkebbySMSGateway(t *testing.T) {
	endpoint := "https://api.skebby.it/API/v1.0/REST/sms"
	recipient := []string{"347....."}
	body := map[string]interface{}{}
	body["returnCredits"] = true
	body["returnCredits"] = true
	body["message"] = "Hello"
	body["message_type"] = "SI"
	body["recipient"] = recipient

	client := httpclient.NewHttpClient()
	client.AddHeader("Content-Type", "application/json")
	client.AddHeader("user_key", "1234")
	client.AddHeader("Session_key", "aDEAQER")
	response, err := client.Post(endpoint, body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(string(response.Body))
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func useResponseBody(body []byte) {
	fmt.Println(string(body))
}
