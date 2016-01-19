package skynology

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/skynology/go-crypto"
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

// 调用自定义函数
func (app *App) Func(name string, data interface{}) (map[string]interface{}, *APIError) {
	_url := fmt.Sprintf("%s/functions/%s", app.baseURL, name)
	result, err := app.sendPostRequest(_url, data)
	return result, err
}

// 调指定url
// url 不包含通用部分. 如 http://skynology.com/api/1.0/files/fetch, 只传入 'files/fetch' 即可
func (app *App) Call(url string, method string, data interface{}) (result map[string]interface{}, err *APIError) {
	_url := fmt.Sprintf("%s/%s", app.baseURL, url)
	method = strings.ToUpper(method)

	switch method {
	case "GET":
		result, err = app.sendGetRequest(_url)
	case "POST":
		result, err = app.sendPostRequest(_url, data)
	case "PUT":
		result, err = app.sendPutRequest(_url, data)
	case "DELETE":
		result, err = app.sendDeleteRequest(_url, data)
	}

	return
}

// 设置处理事件
// 可用在私有项目来避免来回http调用
func (app *App) SetRequestHandler(handler Handler) error {
	app.handler = handler
	return nil
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

	result := fmt.Sprintf("%v,%s", now, crypto.GetMD5(signStr))
	if app.MasterKey != "" {
		result += ",master"
	}

	return result, nil
}

func (app *App) sendGetRequest(url string) (map[string]interface{}, *APIError) {
	return app.sendRequest("GET", url, nil)
}
func (app *App) sendDeleteRequest(url string, data interface{}) (map[string]interface{}, *APIError) {
	return app.sendRequest("DELETE", url, nil)
}

func (app *App) sendPostRequest(url string, data interface{}) (map[string]interface{}, *APIError) {
	return app.sendRequest("POST", url, data)
}

func (app *App) sendPutRequest(url string, data interface{}) (map[string]interface{}, *APIError) {
	return app.sendRequest("PUT", url, data)
}

func (app *App) sendRequest(method string, url string, data interface{}) (map[string]interface{}, *APIError) {
	sign, err := app.getRequestSign()
	if err != nil {
		return nil, &APIError{Code: -1, Error: fmt.Sprintf("marshal json data error:%v", err.Error())}
	}
	params := HandlerRequestParams{}
	params.AppId = app.ApplicationId
	params.AppKey = app.ApplicationKey
	params.MasterKey = app.MasterKey
	params.RequestSign = sign
	params.Method = method
	params.URL = url
	params.SessionToken = app.SessionToken
	params.Data = data

	return app.handler.SendRequest(params)
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

// -----------

func GetInt(v interface{}) int {
	switch reply := v.(type) {
	case int:
		return reply
	case int64:
		x := int(reply)
		if int64(x) != reply {
			return 0
		}
		return x
	case float64:
		x := int(reply)
		if float64(x) != reply {
			return 0
		}
		return x
	case string:
		n, _ := strconv.ParseInt(reply, 10, 0)
		return int(n)
	case []byte:
		n, _ := strconv.ParseInt(string(reply), 10, 0)
		return int(n)
	case nil:
		return 0
	default:
		return 0
	}
	return 0
}

func GetInt64(v interface{}) int64 {
	switch reply := v.(type) {
	case int64:
		return reply
	case int:
		return int64(reply)
	case float64:
		x := int64(reply)
		if float64(x) != reply {
			return 0
		}
		return x
	case string:
		n, _ := strconv.ParseInt(reply, 10, 64)
		return n
	case []byte:
		n, _ := strconv.ParseInt(string(reply), 10, 64)
		return n
	case nil:
		return 0
	default:
		return 0
	}
	return 0
}

func GetFloat64(v interface{}) float64 {
	switch reply := v.(type) {
	case int:
		return float64(reply)
	case int64:
		return float64(reply)
	case string:
		n, _ := strconv.ParseFloat(reply, 64)
		return n
	case []byte:
		n, _ := strconv.ParseFloat(string(reply), 64)
		return n
	case nil:
		return 0
	default:
		return 0
	}
	return 0
}
func GetString(v interface{}) string {
	switch reply := v.(type) {
	case int, int64, float64:
		return fmt.Sprintf("%v", reply)
	case string:
		return reply
	case []byte:
		return string(reply)
	case nil:
		return ""
	default:
		return ""
	}
	return ""
}
func GetBool(v interface{}) bool {

	if v == true {
		return true
	}

	if v == false {
		return true
	}

	switch reply := v.(type) {
	case int, int64, float64:
		return reply != 0
	case string:
		n, _ := strconv.ParseBool(reply)
		return n
	case []byte:
		n, _ := strconv.ParseBool(string(reply))
		return n
	case nil:
		return false
	default:
		return false
	}
	return false
}
