package skynology

import (
	"fmt"
	"strings"
	"time"
)

// get field value
func (obj *Object) Get(field string) interface{} {
	return obj.data[field]
}

// get []interface{}, if field is empty , return []interface{}
func (obj *Object) GetArray(field string) (result []interface{}) {
	v := obj.Get(field)
	if ret, ok := v.([]interface{}); ok {
		result = ret
	}
	return
}

func (obj *Object) GetTime(field string) (result time.Time) {
	v := obj.GetString(field)
	if v == "" {
		return
	}
	layout := time.RFC3339Nano
	if !strings.Contains(v, ".") {
		layout = time.RFC3339
	}

	if t, err := time.Parse(layout, v); err == nil {
		return t
	}
	return
}

func (obj *Object) GetInt(field string) int {
	v := obj.Get(field)
	return GetInt(v)
}

func (obj *Object) GetInt64(field string) int64 {
	v := obj.Get(field)
	return GetInt64(v)
}

func (obj *Object) GetFloat64(field string) float64 {
	v := obj.Get(field)
	return GetFloat64(v)
}
func (obj *Object) GetString(field string) string {
	v := obj.Get(field)
	return GetString(v)
}
func (obj *Object) GetBool(field string) bool {
	v := obj.Get(field)
	return GetBool(v)
}

func (obj *Object) GetMap(field string) (result map[string]interface{}) {
	if m, ok := obj.Get(field).(map[string]interface{}); ok {
		return m
	}

	return
}

// set field value
func (obj *Object) Set(field string, value interface{}) *Object {
	obj.changedData[field] = value
	return obj
}

// 设置多个
func (obj *Object) SetMulti(data map[string]interface{}) *Object {
	for k, v := range data {
		obj.Set(k, v)
	}

	return obj
}

// 返回查到的数据
// 不包括修改的内容
func (obj *Object) Map() map[string]interface{} {
	return obj.data
}

// increment the give amount
func (obj *Object) Increment(field string) *Object {
	return obj.IncrementWithAmount(field, 1)
}

// increment the give amount
func (obj *Object) IncrementWithAmount(field string, amount int) *Object {
	obj.changedData[field] = map[string]interface{}{"__op": "Increment", "amount": amount}
	return obj
}

// add a value to the end of the array field
func (obj *Object) AddValueToArray(field string, value interface{}) *Object {
	obj.changedData[field] = map[string]interface{}{"__op": "Add", "objects": []interface{}{value}}
	return obj
}

// add a value to the array field, only if it is not already present in the array
func (obj *Object) AddUniqueValueToArray(field string, value interface{}) *Object {
	obj.changedData[field] = map[string]interface{}{"__op": "AddUnique", "objects": []interface{}{value}}
	return obj
}

// remove value from array field
func (obj *Object) RemoveValueFromArray(field string, value interface{}) *Object {
	obj.changedData[field] = map[string]interface{}{"__op": "Remove", "objects": []interface{}{value}}
	return obj
}

// add a values from given list to the field
func (obj *Object) AddValueToArrayFromList(field string, value []interface{}) *Object {
	obj.changedData[field] = map[string]interface{}{"__op": "Add", "objects": value}
	return obj
}

// add a values from given list to the field, only if it is not already present in the array
func (obj *Object) AddUniqueValueToArrayFromList(field string, value []interface{}) *Object {
	obj.changedData[field] = map[string]interface{}{"__op": "AddUnique", "objects": value}
	return obj
}

// remove value from array field
func (obj *Object) RemoveValueFromArrayFromList(field string, value []interface{}) *Object {
	obj.changedData[field] = map[string]interface{}{"__op": "Remove", "objects": value}
	return obj
}

// remove object from array
func (obj *Object) RemoveObjectFromArray(field string, query interface{}) *Object {
	obj.changedData[field] = map[string]interface{}{"__op": "RemoveObject", "query": query}
	return obj
}

// update element in array
func (obj *Object) UpdateObjectInArray(query map[string]interface{}, data map[string]interface{}) *Object {
	// just allow update element,
	// so clear all another updates
	obj.changedData = map[string]interface{}{}
	obj.changedData["query"] = query
	obj.changedData["data"] = data
	obj.addtionalURL = "/array"
	return obj
}

// 设置指定用户的读写权限
// 若以往已经有此用户的ACL信息, 将会被覆盖
func (obj *Object) SetReadWriteAccessByUserId(userId string, read bool, write bool) *Object {
	obj.setAccessControl(userId, AccessControlTypeRead, read)
	obj.setAccessControl(userId, AccessControlTypeWrite, write)
	return obj
}

func (obj *Object) SetReadAccessByUserId(userId string, access bool) *Object {
	obj.setAccessControl(userId, AccessControlTypeRead, access)
	return obj
}
func (obj *Object) SetWriteAccessByUserId(userId string, access bool) *Object {
	obj.setAccessControl(userId, AccessControlTypeWrite, access)
	return obj
}

func (obj *Object) SetACL(acl ACL) *Object {
	obj.changedData["ACL"] = acl
	return obj
}

// 用指定角色名来设置读写权限
// 若以往已经有此数据的ACL信息, 将会被覆盖
func (obj *Object) SetReadWriteAccessRoleName(roleName string, read bool, write bool) *Object {
	roleName = "role:" + roleName
	obj.setAccessControl(roleName, AccessControlTypeRead, read)
	obj.setAccessControl(roleName, AccessControlTypeWrite, write)
	return obj
}
func (obj *Object) SetReadAccessByRoleName(roleName string, access bool) *Object {
	roleName = "role:" + roleName
	obj.setAccessControl(roleName, AccessControlTypeRead, access)
	return obj
}
func (obj *Object) SetWriteAccessByRoleName(roleName string, access bool) *Object {
	roleName = "role:" + roleName
	obj.setAccessControl(roleName, AccessControlTypeWrite, access)
	return obj
}
func (obj *Object) setAccessControl(key string, typ AccessControlType, value bool) {
	var acl ACL
	var item ACLItem
	acl, ok := obj.changedData["ACL"].(ACL)
	if !ok {
		acl = ACL{}
	}
	item, ok = acl[key]
	if !ok {
		item = ACLItem{}
	}
	if typ == AccessControlTypeRead {
		item.Read = value
	} else if typ == AccessControlTypeWrite {
		item.Write = value
	}
	acl[key] = item
	obj.changedData["ACL"] = acl
}

// 检查指定用户和角色是否对当前对象有权限
func (obj *Object) CheckACL(userId string, roles []string, typ string) bool {
	acl, err := NewACL(obj.GetMap("ACL"))
	if err != nil {
		return false
	}

	typ = strings.ToLower(typ)

	for k, v := range acl {
		if k == userId {
			if typ == "read" && v.Read {
				return true
			}
			if typ == "write" && v.Write {
				return true
			}
		}
		for _, role := range roles {
			if k == ("role:" + role) {
				if typ == "read" && v.Read {
					return true
				}
				if typ == "write" && v.Write {
					return true
				}
			}
		}
	}
	return false
}

func (obj *Object) Save() (bool, *APIError) {
	var m map[string]interface{}
	var err *APIError

	url := fmt.Sprintf("%s/resources/%s", obj.app.baseURL, obj.ResourceName)
	if obj.ObjectId != "" {
		url += "/" + obj.ObjectId
	}
	if obj.addtionalURL != "" {
		url += obj.addtionalURL
	}

	if obj.ObjectId != "" {
		m, err = obj.app.sendPutRequest(url, obj.changedData)
	} else {
		m, err = obj.app.sendPostRequest(url, obj.changedData)
	}

	fmt.Println("save object:", m)

	if err != nil {
		return false, err
	}

	obj.initData(m)

	return true, nil
}

func (obj *Object) Delete() (bool, *APIError) {
	url := fmt.Sprintf("%s/resources/%s/%s", obj.app.baseURL, obj.ResourceName, obj.ObjectId)
	_, err := obj.app.sendDeleteRequest(url, nil)
	if err != nil {
		return false, err
	}

	obj.clear()

	return true, nil

}

func (obj *Object) initData(data map[string]interface{}) {
	obj.changedData = make(map[string]interface{})
	obj.data = data

	if id, ok := data["objectId"]; ok {
		obj.ObjectId = id.(string)
	}

	if acl, ok := data["ACL"].(map[string]interface{}); ok {
		if nacl, err := NewACL(acl); err == nil {
			obj.ACL = nacl
		}
	}

	if cd, ok := data["createdAt"]; ok {
		if t, err := time.Parse(time.RFC3339Nano, cd.(string)); err == nil {
			obj.CreatedAt = t
		}
	}

	if ud, ok := data["createdAt"]; ok {
		if t, err := time.Parse(time.RFC3339Nano, ud.(string)); err == nil {
			obj.UpdatedAt = t
		}
	}
}

func (obj *Object) clear() {
	obj.changedData = make(map[string]interface{})
	obj.data = make(map[string]interface{})
	obj.ObjectId = ""
	obj.CreatedAt = time.Time{}
	obj.UpdatedAt = time.Time{}
	obj.ACL = nil
}
