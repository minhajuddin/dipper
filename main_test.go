package dipper

import "testing"

type DB string
type User string

type HomeController struct {
	DB
	CurrentUser User
}

func TestRegister(t *testing.T) {
	Register(DB(""), func() interface{} {
		return DB("")
	})
	if len(deps) != 1 {
		t.Errorf("got %v but expected %v", len(deps), 1)
	}
}

func TestInject(t *testing.T) {

	Register(DB(""), func() interface{} {
		return DB("Awesome")
	})

	var db DB
	err, ok := Inject(&db)

	if err != nil || !ok {
		t.Error(err, ok)
	}

	if db != DB("Awesome") {
		t.Errorf("unable to inject %s.", db)
	}
}

func TestInjectStruct(t *testing.T) {

	Register(DB(""), func() interface{} {
		return DB("Awesome")
	})

	var hc HomeController
	err, ok := Inject(&hc)

	if err != nil || !ok {
		t.Error(err, ok)
	}

	if hc.DB != DB("Awesome") {
		t.Errorf("unable to inject %s.", hc.DB)
	}
}

func TestInjectNestedStruct(t *testing.T) {

	Register(DB(""), func() interface{} {
		return DB("Awesome")
	})

	Register(User(""), func() interface{} {
		return User("Zainab")
	})

	type Foo struct {
		Ctrl HomeController
	}

	var f Foo
	err, ok := Inject(&f)

	if err != nil || !ok {
		t.Error(err, ok)
	}

	if f.Ctrl.DB != DB("Awesome") || f.Ctrl.CurrentUser != User("Zainab") {
		t.Errorf("unable to inject %#v.", f)
	}
}
