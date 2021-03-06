package jio

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

func TestArraySchema_SetPriority(t *testing.T) {
	for _, priority := range []int{-1, 0, 100} {
		if priority != Array().SetPriority(priority).Priority() {
			t.Error("set priority failed")
		}
	}
}

func TestArraySchema_TransformAndPrependTransform(t *testing.T) {
	schema := Array().Transform(func(ctx *Context) {
		ctx.Abort(errors.New("2"))
	}).Transform(func(ctx *Context) {
		ctx.Abort(errors.New("3"))
	}).PrependTransform(func(ctx *Context) {
		ctx.Abort(errors.New("1"))
	}).PrependTransform(func(ctx *Context) {
		ctx.Abort(errors.New("0"))
	})
	if len(schema.rules) != 4 {
		t.Error("miss function")
	}
	for i := 0; i < 4; i++ {
		ctx := NewContext(nil)
		schema.rules[i](ctx)
		if ctx.Err.Error() != strconv.Itoa(i) {
			t.Error("sequential error")
		}
	}
}

func TestArraySchema_Required(t *testing.T) {
	schema := Array().Required()
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("should error when no data")
	}
}

func TestArraySchema_Optional(t *testing.T) {
	schema := Array().Optional()
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("should no error")
	}
}

func TestArraySchema_Default(t *testing.T) {
	defaultValue := []int{0, 1, 2, 3}
	schema := Array().Default(defaultValue)
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if reflect.ValueOf(ctx.Value).Len() != 4 {
		t.Error("should set default value")
	}
}

func TestArraySchema_When(t *testing.T) {
	schema := Object().Keys(K{
		"length": String().Required(),
		"list": Array().
			When("length", "2", Array().Length(2)).
			When("length", "3", Array().Length(3)).
			When("length", String(), Array().Min(1)),
	})

	ctx := NewContext(map[string]interface{}{"length": "2", "list": []int{1, 2}})
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("length 2 test failed")
	}

	ctx = NewContext(map[string]interface{}{"length": "3", "list": []int{1, 2}})
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("length 3 test failed")
	}

	ctx = NewContext(map[string]interface{}{"name": "badcase", "age": []int{}})
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("badcase test failed")
	}
}

func TestArraySchema_Check(t *testing.T) {
	schema := Array().Check(func(ctxValue interface{}) error {
		if reflect.ValueOf(ctxValue).Len() != 2 {
			return errors.New("length not equal 2")
		}
		return nil
	})
	ctx := NewContext([]int{1, 2})
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("check should no error")
	}
	ctx = NewContext([]string{"1"})
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("check should error")
	}
	ctx = NewContext("???")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("check should error")
	}
}

func TestArraySchema_Items(t *testing.T) {
	schema := Array().Items(Number().Integer(), String())
	ctx := NewContext([]interface{}{"valid string"})
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("valid string test failed")
	}

	ctx = NewContext([]interface{}{"valid string", 2})
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("valid number test failed")
	}

	ctx = NewContext([]interface{}{"valid string", 3.1})
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("valid decimal test failed")
	}
}

func TestArraySchema_Min(t *testing.T) {
	schema := Array().Min(3)
	ctx := NewContext([]int{0, 1, 2, 3})
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("test min length failed")
	}

	ctx = NewContext([]int{0})
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("test min length should failed")
	}
}

func TestArraySchema_Max(t *testing.T) {
	schema := Array().Max(3)
	ctx := NewContext([]int{0, 1, 2, 3})
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("test max length should failed")
	}

	ctx = NewContext([]int{0})
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("test max length failed")
	}
}

func TestArraySchema_Length(t *testing.T) {
	schema := Array().Max(1)
	ctx := NewContext([]int{0, 1, 2, 3})
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("test length should failed")
	}

	ctx = NewContext([]int{0})
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("test length failed")
	}
}

func TestArraySchema_Validate(t *testing.T) {
	schema := Array()
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("default optional should no error")
	}

	ctx = NewContext("string")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("not array")
	}
}
