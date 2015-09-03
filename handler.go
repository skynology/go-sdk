package skynology

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
)

// 调用http请求时的参数
type HandlerRequestParams struct {
	AppId                string
	AppKey               string
	MasterKey            string
	RequestSign          string
	SessionToken         string
	WeixinId, WeixinType string
	Method               string
	URL                  string
	Data                 interface{}
	Headers              map[string]string
}

// http处理函数,
// 如默认调用http url, 但在私有部署中, 直接调用对应方法, 不走http
type Handler interface {
	SendRequest(params HandlerRequestParams) (map[string]interface{}, *APIError)
}

// 默认处理函数
type DefaultHandler struct{}

func NewDefaultHandler() DefaultHandler {
	return DefaultHandler{}
}
func (d *DefaultHandler) getHttpRequest(params HandlerRequestParams) (*http.Request, error) {
	var body io.Reader
	var length int64
	if params.Data != nil {
		b, err := json.Marshal(params.Data)
		if err != nil {
			return nil, fmt.Errorf("marshal json data error:%v", err.Error())
		}
		body = bytes.NewReader(b)
		length = int64(len(b))
	}

	request, err := http.NewRequest(params.Method, params.URL, body)
	if err != nil {
		return request, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", fmt.Sprintf("Skynology-Golang/%v (%v;%v;)", SDK_VERSION, runtime.GOOS, runtime.GOARCH))

	request.Header.Add(X_CLIENT_VERSION_HEADER, fmt.Sprintf("go-%v", SDK_VERSION))
	request.Header.Add(X_APPLICATION_ID_HEADER, params.AppId)
	request.Header.Add(X_REQUEST_SIGN_HEADER, params.RequestSign)
	if params.SessionToken != "" {
		request.Header.Add(X_SESSION_TOKEN_HEADER, params.SessionToken)
	}
	if params.WeixinId != "" {
		request.Header.Add(X_WEIXIN_ID_HEADER, params.WeixinId)
	}
	if params.WeixinType != "" {
		request.Header.Add(X_WEIXIN_TYPE_HEADER, params.WeixinType)
	}
	request.ContentLength = length

	return request, nil
}

func (d DefaultHandler) SendRequest(params HandlerRequestParams) (map[string]interface{}, *APIError) {
	var apiError APIError
	var m map[string]interface{}

	//fmt.Println("headers:", req.Header)
	//fmt.Println("url:", req.URL)
	req, err := d.getHttpRequest(params)
	if err != nil {
		return nil, &APIError{Code: -1, Error: err.Error()}
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return m, &APIError{Code: -1, Error: fmt.Sprintf("cannot reach skynology server. %v", err.Error())}
	}

	defer response.Body.Close()

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, response.Body)

	//fmt.Println("response is:", string(buf.Bytes()))

	if err != nil {
		return m, &APIError{Code: -1, Error: fmt.Sprintf("cannot read skynology response. %v", err.Error())}
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		err = json.Unmarshal(buf.Bytes(), &m)
		if err != nil {
			return m, &APIError{Code: -1, Error: fmt.Sprintf(" parse response data to json(done). %v", err.Error())}
		}
	} else {
		err = json.Unmarshal(buf.Bytes(), &apiError)
		if err != nil {
			return m, &APIError{Code: -1, Error: fmt.Sprintf("parse response data to json(failed). %v", err.Error())}
		}
		return m, &apiError
	}

	return m, nil
}
