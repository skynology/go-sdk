package skynology

import "fmt"

func (app *App) GetWeixin(url string) (result map[string]interface{}, err *APIError) {
	_url := fmt.Sprintf("%s/weixin/%s", app.baseURL, url)

	result, err = app.sendGetRequest(_url)
	return
}

// 调用前需设置app的 'SetWeixinParams' 方法
// url 无需传入通用部分. 比如创建部门时, 只传"department"即可
// SDK会自动生成完整url
func (app *App) PostWeixin(url string, data interface{}) (result map[string]interface{}, err *APIError) {
	_url := fmt.Sprintf("%s/weixin/%s", app.baseURL, url)
	fmt.Println("send weixin url:", _url)
	result, err = app.sendPostRequest(_url, data)

	return
}

//
func (app *App) PutWeixin(url string, data interface{}) (result map[string]interface{}, err *APIError) {
	_url := fmt.Sprintf("%s/weixin/%s", app.baseURL, url)
	fmt.Println("send weixin url:", _url)
	result, err = app.PutWeixin(_url, data)

	return
}

// 调用前需设置app的 'SetWeixinParams' 方法
func (app *App) DeleteWeixin(url string, data interface{}) (result map[string]interface{}, err *APIError) {
	_url := fmt.Sprintf("%s/weixin/%s", app.baseURL, url)
	fmt.Println("send weixin url:", _url)
	result, err = app.DeleteWeixin(_url, data)

	return
}
