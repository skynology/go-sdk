package skynology

import (
	"fmt"
	"time"
)

// get field value
func (obj *Object) Get(field string) interface{} {
	return obj.data[field]
}

// set field value
func (obj *Object) Set(field string, value interface{}) *Object {
	obj.changedData[field] = value
	return obj
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

func (obj *Object) SetReadAccessByUserId(userId string, access bool) *Object {
	obj.setAccessControl(userId, "read", access)
	return obj
}
func (obj *Object) SetWriteAccessByUserId(userId string, access bool) *Object {
	obj.setAccessControl(userId, "write", access)
	return obj
}

func (obj *Object) SetReadAccessByRoleName(roleName string, access bool) *Object {
	roleName = "role:" + roleName
	obj.setAccessControl(roleName, "read", access)
	return obj
}
func (obj *Object) SetWriteAccessByRoleName(roleName string, access bool) *Object {
	roleName = "role:" + roleName
	obj.setAccessControl(roleName, "write", access)
	return obj
}
func (obj *Object) setAccessControl(key string, typ string, value bool) {
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
	if typ == "read" {
		item.Read = value
	} else if typ == "write" {
		item.Write = value
	}
	acl[key] = item
	obj.changedData["ACL"] = acl
}
func (obj *Object) Save() (bool, *APIError) {
	var m map[string]interface{}
	var err *APIError

	if obj.ObjectId != "" {
		url := fmt.Sprintf("%s/resources/%s/%s", obj.app.baseURL, obj.ResourceName, obj.ObjectId)
		m, err = obj.app.sendPutRequest(url, obj.changedData)
	} else {
		url := fmt.Sprintf("%s/resources/%s", obj.app.baseURL, obj.ResourceName)
		m, err = obj.app.sendPostRequest(url, obj.changedData)
	}

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
