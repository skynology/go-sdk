package skynology

import (
	"fmt"
	"time"
)

// get field value
func (obj *Object) Get(field string) interface{} {
	return obj.data[field]
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
	obj.changedData[field] = map[string]interface{}{"__op": "AddToUnique", "objects": []interface{}{value}}
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
	obj.changedData[field] = map[string]interface{}{"__op": "AddToUnique", "objects": value}
	return obj
}

// remove value from array field
func (obj *Object) RemoveValueFromArrayFromList(field string, value []interface{}) *Object {
	obj.changedData[field] = map[string]interface{}{"__op": "AddToUnique", "objects": value}
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
		nacl, err := NewACL(acl)
		if err != nil {
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
