package skynology

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/williambao/gocrypto"
)

func NewACL(data map[string]interface{}) (ACL, error) {
	acl := ACL{}
	for user, item := range data {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return acl, errors.New("invalid acl format")
		}

		newItem := ACLItem{}
		for k, v := range itemMap {
			if k != "read" && k != "write" {
				return acl, errors.New("invalid acl format")
			}
			bv, ok := v.(bool)
			if !ok {
				return acl, errors.New("invalid acl format")
			}
			if k == "read" {
				newItem.Read = bv
			} else if k == "write" {
				newItem.Write = bv
			}
		}
		acl[user] = newItem
	}

	return acl, nil
}

func (app *App) getRequestSign() (string, error) {
	if app.ApplicationId == "" || app.ApplicationKey == "" && app.MasterKey == "" {
		return "", errors.New("please set `APPLICATION_ID` and `APPLICATION_KEY`")
	}

	now := time.Now().UTC().Unix()
	signStr := fmt.Sprintf("%v%s", now, app.ApplicationKey)

	if app.MasterKey != "" {
		signStr = fmt.Sprintf("%v%s", now, app.MasterKey)
	}

	fmt.Println(signStr)

	result := fmt.Sprintf("%v,%s", now, crypto.GetMD5(signStr))
	if app.MasterKey != "" {
		result += ",master"
	}

	return result, nil
}

func (app *App) getHttpRequest(method string, url string, body io.Reader) (*http.Request, error) {
	sign, err := app.getRequestSign()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return request, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", fmt.Sprintf("Skynology-Golang/%v (%v;%v;)", SDK_VERSION, runtime.GOOS, runtime.GOARCH))

	request.Header.Add(X_CLIENT_VERSION_HEADER, fmt.Sprintf("go-%v", SDK_VERSION))
	request.Header.Add(X_APPLICATION_ID_HEADER, app.ApplicationId)
	request.Header.Add(X_REQUEST_SIGN_HEADER, sign)
	if app.SessionToken != "" {
		request.Header.Add(X_SESSION_TOKEN_HEADER, app.SessionToken)
	}

	return request, nil
}

func (app *App) sendGetRequest(url string) (map[string]interface{}, *APIError) {
	req, err := app.getHttpRequest("GET", url, nil)
	if err != nil {
		return nil, &APIError{Code: -1, Error: err.Error()}
	}
	return app.sendRequest(req)
}
func (app *App) sendDeleteRequest(url string, data interface{}) (map[string]interface{}, *APIError) {
	req, err := app.getHttpRequest("DELETE", url, nil)
	if err != nil {
		return nil, &APIError{Code: -1, Error: err.Error()}
	}
	return app.sendRequest(req)
}
func (app *App) sendPostRequest(url string, data interface{}) (map[string]interface{}, *APIError) {

	b, err := json.Marshal(data)
	if err != nil {
		return nil, &APIError{Code: -1, Error: fmt.Sprintf("marshal json data error:%v", err.Error())}
	}

	req, err := app.getHttpRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, &APIError{Code: -1, Error: err.Error()}
	}
	req.ContentLength = int64(len(b))

	return app.sendRequest(req)
}
func (app *App) sendPutRequest(url string, data interface{}) (map[string]interface{}, *APIError) {

	b, err := json.Marshal(data)
	if err != nil {
		return nil, &APIError{Code: -1, Error: fmt.Sprintf("marshal json data error:%v", err.Error())}
	}

	req, err := app.getHttpRequest("PUT", url, bytes.NewReader(b))
	if err != nil {
		return nil, &APIError{Code: -1, Error: err.Error()}
	}
	req.ContentLength = int64(len(b))

	return app.sendRequest(req)
}
func (app *App) sendRequest(req *http.Request) (map[string]interface{}, *APIError) {
	var apiError *APIError
	var m map[string]interface{}

	//fmt.Println("headers:", req.Header)
	//fmt.Println("url:", req.URL)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return m, &APIError{Code: -1, Error: fmt.Sprintf("cannot reach skynology server. %v", err.Error())}
	}

	defer response.Body.Close()

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, response.Body)

	if err != nil {
		return m, &APIError{Code: -1, Error: fmt.Sprintf("cannot read skynology response. %v", err.Error())}
	}

	//fmt.Println("query response data:", string(buf.Bytes()))

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		err = json.Unmarshal(buf.Bytes(), &m)
	} else {
		err = json.Unmarshal(buf.Bytes(), &apiError)
	}

	if err != nil {
		return m, &APIError{Code: -1, Error: fmt.Sprintf("cannot read skynology response. %v", err.Error())}
	}

	return m, nil
}

func (app *App) saveUserToDisk(user *User) error {
	bin, err := json.Marshal(user.data)
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%ssynology_session_%s", app.dataDir, app.ApplicationId)

	err = ioutil.WriteFile(filePath, bin, os.ModePerm)
	if err != nil {
		return err
	}

	user.app.currentUser = user

	return nil
}

func (app *App) getUserFromDisk() (*User, error) {

	var m map[string]interface{}
	filePath := fmt.Sprintf("%ssynology_session_%s", app.dataDir, app.ApplicationId)

	bin, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bin, &m)
	if err != nil {
		return nil, err
	}

	user := app.NewUserWithData(m)
	return user, nil
}

func (app *App) clearUserFromDisk() error {
	filePath := fmt.Sprintf("%ssynology_session_%s", app.dataDir, app.ApplicationId)
	err := ioutil.WriteFile(filePath, []byte(""), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
