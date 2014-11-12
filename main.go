//Package dipper is a small dependency injection library
package dipper

import (
	"fmt"
	"reflect"
	"sync"
)

//Constructor is a function which knows how to construct a type
type Constructor func() interface{}

var deps = make(map[reflect.Type]Constructor)

var mutex = sync.RWMutex{}

//Register registers a type and its constructor function.
//e.g.
//    dipper.Register(sql.DB, func()interface{}{
//      //code to create the db and return it
//      return db;
//    }
func Register(val interface{}, c Constructor) {
	mutex.Lock()
	defer mutex.Unlock()
	deps[reflect.TypeOf(val)] = c
}

//MustInject is similar to inject but panics if
//it finds a dependency and is unable to inject it
func MustInject(val interface{}) bool {
	err, ok := Inject(val)
	if err != nil {
		panic(err)
	}
	return ok
}

//Inject injects the right values into the dependency
//e.g.
//    dipper.Register(sql.DB, func()interface{}{
//      //code to create the db and return it
//      return db;
//    }
//    db := sql.DB{}
//    dipper.Inject(&db)
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
