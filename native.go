package goja

import (
	"bytes"
	"errors"

	"github.com/dop251/goja/unistring"
)

type Property struct {
	Name  string
	Value Value
}

type NativeClass struct {
	*Object
	runtime    *Runtime
	classProto *Object
	className  string
	ctor       func(call FunctionCall) interface{}
	classProps []Property
	funcProps  []Property

	Function *Object

	getStacktrace func(err error) string
}

func (r *Runtime) TryToValue(i interface{}) (Value, error) {
	var result Value
	err := r.vm.try(r.vm.ctx, func() {
		result = r.ToValue(i)
	})
	if err.Error() == "" || err.Error() == "<nil>" {
		return result, nil
	}
	return result, err
}

func (r *Runtime) MakeCustomError(name, msg string) *Object {
	e := r.newError(r.global.Error, msg).(*Object)
	e.self.setOwnStr("name", asciiString(name), false)
	e.self.setOwnStr("customerror", TrueValue(), false)
	return e
}

func (r *Runtime) CreateNativeErrorClass(
	className string,
	ctor func(call FunctionCall) error,
	initStacktrace func(err error, stacktrace string),
	getStacktrace func(err error) string,
	classProps []Property,
	funcProps []Property,
) NativeClass {
	//TODO goja handle getStacktrace
	classProto := r.builtin_new(r.global.Error, []Value{})
	o := classProto.self
	o._putProp("name", asciiString(className), true, false, true)

	for _, prop := range classProps {
		o._putProp(unistring.String(prop.Name), prop.Value, true, false, true)
	}

	var errMsg valueString
	v := r.newNativeFuncConstruct(func(args []Value, proto *Object) *Object {
		obj := r.newBaseObject(proto, className)
		call := FunctionCall{
			ctx:       r.vm.ctx,
			This:      obj.val,
			Arguments: args,
		}

		err := ctor(call)
		ex := &Exception{
			val: r.newError(r.global.ReferenceError, err.Error()),
		}
		stackTrace := bytes.NewBuffer(nil)
		ex.writeShortStack(stackTrace)
		initStacktrace(err, stackTrace.String())
		errMsg = newStringValue(err.Error())
		obj._putProp("message", errMsg, true, false, true)

		g := &_goNativeValue{baseObject: obj, value: err}
		obj.val.self = g
		obj.val.__wrapped = g.value

		return obj.val
	}, unistring.String(className), classProto, 1)

	for _, prop := range funcProps {
		v.self._putProp(unistring.String(prop.Name), prop.Value, true, false, true)
	}
	v.runtime = r

	return NativeClass{Object: v, runtime: r, classProto: classProto, className: className, Function: v, getStacktrace: getStacktrace}
}

func (r *Runtime) CreateNativeError(name string) (Value, func(err error) Value) {
	proto := r.builtin_new(r.global.Error, []Value{})
	o := proto.self
	o._putProp("name", asciiString(name), true, false, true)

	e := r.newNativeFuncConstructProto(r.builtin_Error, unistring.String(name), proto, r.global.Error, 1)

	return e, func(err error) Value {
		return r.MakeCustomError(name, err.Error())
	}
}

func (r *Runtime) CreateNativeClass(
	className string,
	ctor func(call FunctionCall) interface{},
	classProps []Property,
	funcProps []Property,
) NativeClass {
	// needed for instance of
	classProto := r.builtin_new(r.global.Object, []Value{})
	o := classProto.self
	o._putProp("name", asciiString(className), true, false, true)

	for _, prop := range classProps {
		o._putProp(unistring.String(prop.Name), prop.Value, true, false, true)
	}

	ctorImpl := func(call ConstructorCall) *Object {
		fCall := FunctionCall{
			ctx:       r.vm.ctx,
			This:      call.This,
			Arguments: call.Arguments,
		}
		val := ctor(fCall)
		call.This.__wrapped = val
		// add the toString function first so it can be overridden if user wants to do so
		call.This.self._putProp("name", asciiString(className), true, false, true)
		for _, prop := range funcProps {
			call.This.self._putProp(unistring.String(prop.Name), prop.Value, true, false, true)
		}
		for _, prop := range classProps {
			call.This.self._putProp(unistring.String(prop.Name), prop.Value, true, false, true)
		}
		return nil
	}
	responseObject := r.ToValue(ctorImpl).(*Object)
	p := responseObject.Get("prototype")
	pObject := p.(*Object)
	proto := pObject.self

	proto._putProp("name", asciiString(className), true, false, true)
	for _, prop := range classProps {
		proto.setOwnStr(unistring.String(prop.Name), prop.Value, false)
	}
	responseObject.self._putProp("name", asciiString(className), true, false, true)
	for _, prop := range funcProps {
		responseObject.self._putProp(unistring.String(prop.Name), prop.Value, true, false, true)
	}

	for _, prop := range funcProps {
		proto._putProp(unistring.String(prop.Name), prop.Value, true, false, true)
	}

	v := NativeClass{classProto: pObject, className: className, classProps: classProps, funcProps: funcProps, ctor: ctor, Function: responseObject}
	v.runtime = r

	return v
}

type _goNativeValue struct {
	*baseObject
	value interface{}
}

func (n NativeClass) InstanceOf(val interface{}) Value {
	r := n.runtime
	className := n.className
	classProto := n.classProto
	obj, err := r.New(r.newNativeFuncConstruct(func(args []Value, proto *Object) *Object {
		obj := r.newBaseObject(proto, className)
		obj.class = n.className
		g := &_goNativeValue{baseObject: obj, value: val}
		obj.val.self = g
		obj.val.__wrapped = g.value
		return obj.val
	}, unistring.String(className), classProto, 1))
	if err != nil {
		panic(err)
	}
	if err, ok := val.(error); ok {
		if n.getStacktrace != nil {
			stackTrace := n.getStacktrace(err)
			if len(stackTrace) == 0 {
				obj.self.setOwnStr("customerror", TrueValue(), false)
			}
		}
		obj.self._putProp("message", newStringValue(err.Error()), true, false, true)
		// stack trace support
	}
	obj.self._putProp("name", asciiString(n.className), true, false, true)
	for _, prop := range n.funcProps {
		obj.self._putProp(unistring.String(prop.Name), prop.Value, true, false, true)
	}
	return obj
}

// CreateNativeFunction creates a native function that will call the given call function.
// This provides for a way to detail how the function appears to a user within JS
// compared to passing the call in via toValue.
func (r *Runtime) CreateNativeFunction(name, file string, call func(FunctionCall) Value) (Value, error) {
	if call == nil {
		return UndefinedValue(), errors.New("call cannot be nil")
	}

	return r.newNativeFunc(call, nil, unistring.String(name), nil, 1), nil
}

func (r *Runtime) Eval(name, src string, direct, strict bool) (Value, error) {
	this := r.NewObject()

	p, err := r.compile(name, src, strict, true)
	if err != nil {
		panic(err)
	}

	vm := r.vm

	vm.pushCtx()
	vm.prg = p
	vm.pc = 0
	if !direct {
		vm.stash = nil
	}
	vm.sb = vm.sp
	vm.push(this)
	if strict {
		vm.push(valueTrue)
	} else {
		vm.push(valueFalse)
	}
	vm.run()
	vm.popCtx()
	vm.halt = false
	retval := vm.stack[vm.sp-1]
	vm.sp -= 2
	return retval, nil
}
