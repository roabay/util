package resource

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/rs/rest-layer/schema"
	"github.com/stretchr/testify/assert"
)

func TestValidateSelector(t *testing.T) {
	s := schema.Schema{
		Fields: schema.Fields{
			"parent": {
				Schema: &schema.Schema{
					Fields: schema.Fields{"child": {}},
				},
			},
			"simple": schema.Field{},
			"with_params": {
				Params: schema.Params{
					"foo": {
						Validator: schema.Integer{},
					},
				},
			},
		},
	}

	assert.NoError(t, validateSelector([]Field{{Name: "parent", Fields: []Field{{Name: "child"}}}}, s))
	assert.NoError(t, validateSelector([]Field{{Name: "with_params", Params: map[string]interface{}{"foo": 1}}}, s))

	assert.EqualError(t,
		validateSelector([]Field{{Name: "foo"}}, s),
		"foo: unknown field")
	assert.EqualError(t,
		validateSelector([]Field{{Name: "simple", Fields: []Field{{Name: "child"}}}}, s),
		"simple: field as no children")
	assert.EqualError(t,
		validateSelector([]Field{{Name: "parent", Fields: []Field{{Name: "foo"}}}}, s),
		"parent.foo: unknown field")
	assert.EqualError(t,
		validateSelector([]Field{{Name: "simple", Params: map[string]interface{}{"foo": 1}}}, s),
		"simple: params not allowed")
	assert.EqualError(t,
		validateSelector([]Field{{Name: "with_params", Params: map[string]interface{}{"bar": 1}}}, s),
		"with_params: unsupported param name: bar")
	assert.EqualError(t,
		validateSelector([]Field{{Name: "with_params", Params: map[string]interface{}{"foo": "a string"}}}, s),
		"with_params: invalid param `foo' value: not an integer")
}

func TestApplySelector(t *testing.T) {
	s := schema.Schema{
		Fields: schema.Fields{
			"parent": {
				Schema: &schema.Schema{
					Fields: schema.Fields{
						"child": {},
					},
				},
			},
			"simple": schema.Field{},
			"with_params": {
				Params: schema.Params{
					"foo": {Validator: schema.Integer{}},
				},
				Handler: func(ctx context.Context, value interface{}, params map[string]interface{}) (interface{}, error) {
					if val, found := params["foo"]; found {
						if val == -1 {
							return nil, errors.New("some error")
						}
						return fmt.Sprintf("param is %d", val), nil
					}
					return "no param", nil
				},
			},
		},
	}

	// Basic filtering
	ctx := context.Background()
	p, err := applySelector(ctx, []Field{{Name: "parent", Fields: []Field{{Name: "child"}}}}, s,
		map[string]interface{}{"parent": map[string]interface{}{"child": "value"}, "simple": "value"}, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, map[string]interface{}{"parent": map[string]interface{}{"child": "value"}}, p)
	}
	// Alias on both parent and child
	p, err = applySelector(ctx, []Field{{Name: "parent", Alias: "p", Fields: []Field{{Name: "child", Alias: "c"}}}}, s,
		map[string]interface{}{"parent": map[string]interface{}{"child": "value"}}, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, map[string]interface{}{"p": map[string]interface{}{"c": "value"}}, p)
	}
	// Param call with valid value
	p, err = applySelector(ctx, []Field{{Name: "with_params", Params: map[string]interface{}{"foo": 1}}}, s,
		map[string]interface{}{"with_params": "value"}, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, map[string]interface{}{"with_params": "param is 1"}, p)
	}
	// If no param, handler do not call handler
	p, err = applySelector(ctx, []Field{{Name: "with_params"}}, s,
		map[string]interface{}{"with_params": "value"}, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, map[string]interface{}{"with_params": "value"}, p)
	}
	// Param call with valid value rejected by the handler
	p, err = applySelector(ctx, []Field{{Name: "with_params", Params: map[string]interface{}{"foo": -1}}}, s,
		map[string]interface{}{"with_params": "value"}, nil)
	assert.EqualError(t, err, "with_params: some error")
	assert.Nil(t, p)
	// Deep field lookup on a field with no child
	p, err = applySelector(ctx, []Field{{Name: "simple", Fields: []Field{{Name: "child"}}}}, s,
		map[string]interface{}{"simple": "value"}, nil)
	assert.EqualError(t, err, "simple: field as no children")
	assert.Nil(t, p)
	// Deep field lookup on a field with invalid payload (no dict)
	p, err = applySelector(ctx, []Field{{Name: "parent", Fields: []Field{{Name: "child"}}}}, s,
		map[string]interface{}{"parent": "value"}, nil)
	assert.EqualError(t, err, "parent: invalid value: not a dict")
	assert.Nil(t, p)
}
