package mongo

import (
	"testing"

	"regexp"

	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/schema"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

type UnsupportedExpression struct{}

func (u UnsupportedExpression) Match(p map[string]interface{}) bool {
	return false
}

func callGetQuery(q schema.Query) (bson.M, error) {
	l := resource.NewLookup()
	l.AddQuery(q)
	return getQuery(l)
}

func callGetSort(s string, v schema.Validator) []string {
	l := resource.NewLookup()
	l.SetSort(s, v)
	return getSort(l)
}

func TestGetQuery(t *testing.T) {
	var b bson.M
	var err error
	b, err = callGetQuery(schema.Query{schema.Equal{Field: "id", Value: "foo"}})
	assert.NoError(t, err)
	assert.Equal(t, bson.M{"_id": "foo"}, b)
	b, err = callGetQuery(schema.Query{schema.Equal{Field: "f", Value: "foo"}})
	assert.NoError(t, err)
	assert.Equal(t, bson.M{"f": "foo"}, b)
	b, err = callGetQuery(schema.Query{schema.NotEqual{Field: "f", Value: "foo"}})
	assert.NoError(t, err)
	assert.Equal(t, bson.M{"f": bson.M{"$ne": "foo"}}, b)
	b, err = callGetQuery(schema.Query{schema.GreaterThan{Field: "f", Value: 1}})
	assert.NoError(t, err)
	assert.Equal(t, bson.M{"f": bson.M{"$gt": float64(1)}}, b)
	b, err = callGetQuery(schema.Query{schema.GreaterOrEqual{Field: "f", Value: 1}})
	assert.NoError(t, err)
	assert.Equal(t, bson.M{"f": bson.M{"$gte": float64(1)}}, b)
	b, err = callGetQuery(schema.Query{schema.LowerThan{Field: "f", Value: 1}})
	assert.NoError(t, err)
	assert.Equal(t, bson.M{"f": bson.M{"$lt": float64(1)}}, b)
	b, err = callGetQuery(schema.Query{schema.LowerOrEqual{Field: "f", Value: 1}})
	assert.NoError(t, err)
	assert.Equal(t, bson.M{"f": bson.M{"$lte": float64(1)}}, b)
	b, err = callGetQuery(schema.Query{schema.In{Field: "f", Values: []schema.Value{"foo", "bar"}}})
	assert.NoError(t, err)
	assert.Equal(t, bson.M{"f": bson.M{"$in": []interface{}{"foo", "bar"}}}, b)
	b, err = callGetQuery(schema.Query{schema.NotIn{Field: "f", Values: []schema.Value{"foo", "bar"}}})
	assert.NoError(t, err)
	assert.Equal(t, bson.M{"f": bson.M{"$nin": []interface{}{"foo", "bar"}}}, b)
	if v, err := regexp.Compile("fo[o]{1}.+is.+some"); err == nil {
		b, err = callGetQuery(schema.Query{schema.Regex{Field: "f", Value: v}})
	}
	assert.NoError(t, err)
	assert.Equal(t, bson.M{"f": bson.M{"$regex": "fo[o]{1}.+is.+some"}}, b)
	b, err = callGetQuery(schema.Query{schema.And{schema.Equal{Field: "f", Value: "foo"}, schema.Equal{Field: "f", Value: "bar"}}})
	assert.NoError(t, err)
	assert.Equal(t, bson.M{"$and": []bson.M{bson.M{"f": "foo"}, bson.M{"f": "bar"}}}, b)
	b, err = callGetQuery(schema.Query{schema.Or{schema.Equal{Field: "f", Value: "foo"}, schema.Equal{Field: "f", Value: "bar"}}})
	assert.NoError(t, err)
	assert.Equal(t, bson.M{"$or": []bson.M{bson.M{"f": "foo"}, bson.M{"f": "bar"}}}, b)
}

func TestGetQueryInvalid(t *testing.T) {
	var err error
	_, err = callGetQuery(schema.Query{UnsupportedExpression{}})
	assert.Equal(t, resource.ErrNotImplemented, err)
	_, err = callGetQuery(schema.Query{schema.And{UnsupportedExpression{}}})
	assert.Equal(t, resource.ErrNotImplemented, err)
	_, err = callGetQuery(schema.Query{schema.Or{UnsupportedExpression{}}})
	assert.Equal(t, resource.ErrNotImplemented, err)
}

func TestGetSort(t *testing.T) {
	var s []string
	v := schema.Schema{Fields: schema.Fields{"id": schema.IDField, "f": {Sortable: true}}}
	s = callGetSort("", v)
	assert.Equal(t, []string{"_id"}, s)
	s = callGetSort("id", v)
	assert.Equal(t, []string{"_id"}, s)
	s = callGetSort("f", v)
	assert.Equal(t, []string{"f"}, s)
	s = callGetSort("-f", v)
	assert.Equal(t, []string{"-f"}, s)
	s = callGetSort("f,-f", v)
	assert.Equal(t, []string{"f", "-f"}, s)
}
