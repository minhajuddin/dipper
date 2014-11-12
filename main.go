package dipper

import (
	"fmt"
	"reflect"
	"sync"
)

type Constructor func() interface{}

var deps = make(map[reflect.Type]Constructor)

var mutex = sync.RWMutex{}

//Register
func Register(val interface{}, c Constructor) {
	mutex.Lock()
	defer mutex.Unlock()
	deps[reflect.TypeOf(val)] = c
}

//Inject
func Inject(val interface{}) (error, bool) {
	mutex.RLock()
	defer mutex.RUnlock()

	v := reflect.ValueOf(val)

	//if value is not settable we can't do anything
	if v.Kind() != reflect.Ptr || !v.Elem().CanSet() {
		return fmt.Errorf("input val is not settable type: '%v', val: '%v'", v.Kind(), v.Interface()), true
	}

	t := v.Elem().Type()
	ctor, ok := deps[t]

	//no constructor for this type
	//and it is not a struct, so just return
	if !ok && t.Kind() != reflect.Struct {
		return nil, false
	}

	//we have a constructor for this
	if ok {
		cv := ctor()
		v.Elem().Set(reflect.ValueOf(cv))
		return nil, true
	}

	//this is a struct, let's fill the dependencies
	//of all its exported types
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		cv := reflect.New(f.Type)
		err, ok := Inject(cv.Interface())
		//an error occured when trying to inject a field value
		if err != nil {
			return err, true
		}
		//value was injected
		if ok {
			v.Elem().Field(i).Set(cv.Elem())
		}
	}

	return nil, true
}
