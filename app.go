package skynology

import "fmt"

// create a new Skynology sdk instance
func NewApp(appId, appKey string) *App {
	return &App{
		ApplicationId:  appId,
		ApplicationKey: appKey,
		baseURL:        "https://skynology.com/api/1.0",
		dataDir:        "./",
		weixinParams:   new(weixinParams),
	}
}

// create a new Skynology sdk instance and set app id & master key
func NewAppWithMasterKey(appId, masterKey string) *App {
	return &App{
		ApplicationId: appId,
		MasterKey:     masterKey,
		baseURL:       "https://skynology.com/api/1.0",
		dataDir:       "./",
		weixinParams:  new(weixinParams),
	}
}

// create Object
// create new empty object
func (app *App) NewObject(resourceName string) *Object {
	return &Object{
		app:          app,
		ResourceName: resourceName,
		data:         make(map[string]interface{}),
		changedData:  make(map[string]interface{}),
		baseURL:      fmt.Sprintf("%s/resources/%s", app.baseURL, resourceName),
	}
}

// create object with objectId
func (app *App) NewObjectWithId(resourceName string, objectId string) *Object {
	return &Object{
		app:          app,
		ResourceName: resourceName,
		ObjectId:     objectId,
		data:         make(map[string]interface{}),
		changedData:  make(map[string]interface{}),
		baseURL:      fmt.Sprintf("%s/resources/%s", app.baseURL, resourceName),
	}
}

// create object with initialize data
func (app *App) NewObjectWithData(resourceName string, data map[string]interface{}) *Object {
	obj := &Object{
		app:          app,
		ResourceName: resourceName,
		baseURL:      fmt.Sprintf("%s/resources/%s", app.baseURL, resourceName),
	}

	obj.initData(data)

	return obj
}

// create Query
func (app *App) NewQuery(resourceName string) *Query {
	return &Query{
		app:          app,
		ResourceName: resourceName,
		_take:        20,
		where:        make(map[string]interface{}),
	}
}

// create new user
func (app *App) NewUser() *User {
	user := &User{
		Object: Object{
			app:          app,
			ResourceName: "_User",
			data:         make(map[string]interface{}),
			changedData:  make(map[string]interface{}),
			baseURL:      fmt.Sprintf("%s/users", app.baseURL),
		},
	}

	return user
}

func (app *App) NewUserWithId(objectId string) *User {
	user := &User{
		Object: Object{
			app:          app,
			ResourceName: "_User",
			data:         make(map[string]interface{}),
			changedData:  make(map[string]interface{}),
			baseURL:      fmt.Sprintf("%s/users", app.baseURL),
		},
	}
	user.ObjectId = objectId

	return user
}

func (app *App) NewUserWithData(data map[string]interface{}) *User {
	user := &User{
		Object: Object{
			app:          app,
			ResourceName: "_User",
			data:         make(map[string]interface{}),
			changedData:  make(map[string]interface{}),
			baseURL:      fmt.Sprintf("%s/users", app.baseURL),
		},
	}

	user.initData(data)

	return user
}

// get logined user
func (app *App) CurrentUser() *User {
	user, err := app.getUserFromDisk()
	if err != nil {
		return nil
	}
	return user
}

func (app *App) SetBaseURL(url string) {
	app.baseURL = url
}

// 设置微信配置
// id : 在管理后台绑定完微信公众号后由系统产生的Id
// typ: 公众号类型, 如corp, mp等.
func (app *App) InitWeixin(id string, typ string) {
	app.weixinParams = &weixinParams{Id: id, Type: typ}
}
