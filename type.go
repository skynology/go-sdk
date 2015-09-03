package skynology

import (
	"fmt"
	"time"
)

// API HTTP method
// Can be GET, POST, PUT or DELETE
type Method string
type AccessControlType byte

// Object ACL attribute type
type ACLItem struct {
	Read  bool `json:"read"`
	Write bool `json:"write"`
}
type ACL map[string]ACLItem

const (
	AccessControlTypeRead AccessControlType = iota
	AccessControlTypeWrite
)

// API params.
//
// For general uses, just use Params as a ordinary map.
//
// For advanced uses, use MakeParams to create Params from any struct.
type Params map[string]interface{}

// Skynology API call result.
type Result map[string]interface{}

// 微信配置
type weixinParams struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

func (p *weixinParams) IsValid() bool {
	if p.Id == "" || p.Type == "" {
		return false
	}
	return true
}

// Skynology API error.
type APIError struct {
	Code        int    `json:"code"`
	EnError     string `json:"error_en"`
	Error       string `json:"error"`
	Description string `json:"description,omitempty"`
}

func (a *APIError) String() string {
	return fmt.Sprintf("code:%v, error:%s, description:%s", a.Code, a.Error, a.Description)
}

// Skynology GO SDK app
type App struct {
	ApplicationId  string
	ApplicationKey string
	MasterKey      string
	SessionToken   string
	baseURL        string
	dataDir        string
	currentUser    *User
	weixinParams   *weixinParams
	handler        Handler
}

// query function
type Query struct {
	app          *App
	ResourceName string
	_skip        int
	_take        int
	_count       bool
	where        map[string]interface{}
	order        []string
	field        []string
	include      []string
}

// object type
type Object struct {
	app          *App
	ResourceName string
	ObjectId     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ACL          ACL
	data         map[string]interface{}
	changedData  map[string]interface{}
	baseURL      string

	// 特殊API需额外URL时, 如更新数组字段内对象时, 设置此
	addtionalURL string
}

type User struct {
	Object
	UserName string
	Password string
	Email    string
	Phone    string
}

type File struct {
	Object
	Key    string
	Url    string
	Bucket string
}

// geoJSON
type CoordType float64
type Coordinate [2]CoordType
type Coordinates []Coordinate
type MultiLine []Coordinates
