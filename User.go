package skynology

import (
	"fmt"
	"log"
)

func (user *User) Register() (bool, *APIError) {
	var m map[string]interface{}
	var err *APIError
	url := user.baseURL
	m, err = user.app.sendPostRequest(url, user.changedData)
	if err != nil {
		return false, err
	}

	user.initData(m)

	return true, nil
}

// 重设置密码
// 调用此方法前需已经登录
func (user *User) ResetPassword(oldPassword string, newPassword string) (bool, *APIError) {
	data := map[string]interface{}{
		"old_password": oldPassword,
		"new_password": newPassword,
	}

	url := fmt.Sprintf("%s/%s/resetPassword", user.baseURL, user.ObjectId)

	_, err := user.app.sendPostRequest(url, data)
	if err != nil {
		return false, err
	}

	return true, nil
}

// 使用用户名和密码登录
func (app *App) LoginWithUserName(username, password string) (*User, *APIError) {
	data := map[string]interface{}{
		"username": username,
		"password": password,
	}
	return app.login(data)
}

// 使用手机号码和密码登录
func (app *App) LoginWithPhone(phone, password string) (*User, *APIError) {
	data := map[string]interface{}{
		"phone":    phone,
		"password": password,
	}
	return app.login(data)
}

// 使用邮箱和密码登录
func (app *App) LoginWithEmail(email, password string) (*User, *APIError) {
	data := map[string]interface{}{
		"email":    email,
		"password": password,
	}
	return app.login(data)
}

func (user *User) Logout() (bool, *APIError) {
	url := fmt.Sprintf("%s/logout", user.app.baseURL)
	data := map[string]interface{}{
		"objectId": user.ObjectId,
	}

	_, err := user.app.sendPostRequest(url, data)
	if err != nil {
		return false, &APIError{Code: -1, Error: err.Error}
	}

	err2 := user.app.clearUserFromDisk()
	if err2 != nil {
		return false, &APIError{Code: -1, Error: err2.Error()}
	}

	return true, nil
}

// 覆盖父struct的方法
func (user *User) initData(data map[string]interface{}) {
	user.Object.initData(data)
	if v, ok := data["username"]; ok {
		user.UserName = v.(string)
	}
	if v, ok := data["email"]; ok {
		user.Email = v.(string)
	}
	if v, ok := data["phone"]; ok {
		user.Phone = v.(string)
	}
}

// 覆盖父struct的方法
func (user *User) clear() {
	user.Object.clear()
	user.UserName = ""
	user.Password = ""
	user.Phone = ""
	user.Email = ""
}

func (app *App) login(data map[string]interface{}) (*User, *APIError) {
	url := fmt.Sprintf("%s/login", app.baseURL)

	m, err := app.sendPostRequest(url, data)
	if err != nil {
		return nil, err
	}

	user := app.NewUserWithData(m)

	// save to local disk
	if err := app.saveUserToDisk(user); err != nil {
		log.Printf("save user to disk failed, error:%v", err)
	}

	return user, nil
}
