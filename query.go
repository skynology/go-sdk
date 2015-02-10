package skynology

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func (query *Query) Equal(field string, value interface{}) *Query {
	query.where[field] = value
	return query
}

func (query *Query) NotEqual(field string, value interface{}) *Query {
	query.where[field] = map[string]interface{}{"$ne": value}
	return query
}

func (query *Query) LessThan(field string, value float64) *Query {
	query.where[field] = map[string]interface{}{"$lt": value}
	return query
}
func (query *Query) LessThanOrEqual(field string, value float64) *Query {
	query.where[field] = map[string]interface{}{"$lte": value}
	return query
}
func (query *Query) GreaterThan(field string, value float64) *Query {
	query.where[field] = map[string]interface{}{"$gt": value}
	return query
}
func (query *Query) GreaterThanOrEqual(field string, value float64) *Query {
	query.where[field] = map[string]interface{}{"$gte": value}
	return query
}

func (query *Query) StartWith(field string, value string) *Query {
	query.where[field] = map[string]interface{}{"$regex": fmt.Sprintf("/^%s/i", value)}
	return query
}
func (query *Query) EndWith(field string, value string) *Query {
	query.where[field] = map[string]interface{}{"$regex": fmt.Sprintf("/%s$/i", value)}
	return query
}
func (query *Query) Contains(field string, value string) *Query {
	query.where[field] = map[string]interface{}{"$regex": fmt.Sprintf("/%s/i", value)}
	return query
}

func (query *Query) Exists(field string, exist bool) *Query {
	query.where[field] = map[string]interface{}{"$exists": exist}
	return query
}

func (query *Query) Count(value bool) *Query {
	query._count = value
	return query
}

func (query *Query) Skip(value int) *Query {
	query._skip = value
	return query
}

func (query *Query) Take(value int) *Query {
	query._take = value
	return query
}

func (query *Query) In(field string, value []interface{}) *Query {
	query.where[field] = map[string]interface{}{"$in": value}
	return query
}
func (query *Query) NotIn(field string, value []interface{}) *Query {
	query.where[field] = map[string]interface{}{"$nin": value}
	return query
}
func (query *Query) MatchAll(field string, value []interface{}) *Query {
	query.where[field] = map[string]interface{}{"$all": value}
	return query
}

func (query *Query) OrderBy(field string) *Query {
	query.order = append(query.order, field)
	return query
}
func (query *Query) OrderByDescending(field string) *Query {
	field = "-" + field
	query.order = append(query.order, field)
	return query
}

func (query *Query) Select(fields ...string) *Query {
	query.field = append(query.field, fields...)
	return query
}

func (query *Query) GetObject(objectId string) (Object, *APIError) {
	var result Object

	url := fmt.Sprintf("%s/resources/%s/%s", query.app.baseURL, query.ResourceName, objectId)
	m, err := query.app.sendGetRequest(url)
	if err != nil {
		return result, err
	}

	result = *query.app.NewObjectWithData(query.ResourceName, m)

	return result, nil
}

// 返回 数据列表， 总数 及出错信息
func (query *Query) Find() ([]Object, int, *APIError) {
	var result []Object
	var count int64 = 0

	url := fmt.Sprintf("%s/resources/%s?%s", query.app.baseURL, query.ResourceName, query.getQueryString())

	m, err := query.app.sendGetRequest(url)
	if err != nil {
		return result, 0, err
	}

	if results, ok := m["results"].([]interface{}); ok {
		for _, c := range results {
			obj := query.app.NewObjectWithData(query.ResourceName, c.(map[string]interface{}))
			result = append(result, *obj)
		}
	}

	if c, ok := m["count"].(float64); ok {
		count = int64(c)
	}

	return result, int(count), nil
}

func (query *Query) getQueryString() string {

	search := "_=_"
	if query._count {
		search += "&count=1"
	}
	if len(query.order) > 0 {
		search += ("&order=" + strings.Join(query.order, ","))
	}
	if len(query.field) > 0 {
		search += ("&select=" + strings.Join(query.field, ","))
	}
	if query._skip > 0 {
		search += "&skip=" + strconv.Itoa(query._skip)
	}

	search += "&take=" + strconv.Itoa(query._take)

	if b, err := json.Marshal(query.where); err == nil {
		where := url.QueryEscape(string(b))
		search += "&where=" + where
	}

	return search
}
