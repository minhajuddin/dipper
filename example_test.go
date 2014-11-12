package dipper_test

import (
	"fmt"

	"github.com/minhajuddin/dipper"
)

type DB string
type User string

type HomeController struct {
	DB
	CurrentUser User
}

func ExampleInject() {
	//one time registration of types and their constructors
	dipper.Register(DB(""), func() interface{} {
		return DB("Awesome")
	})

	dipper.Register(User(""), func() interface{} {
		return User("Zainab")
	})

	hc := HomeController{}

	dipper.Inject(&hc)

	fmt.Printf("%#v", hc)
	// Output:
	// dipper_test.HomeController{DB:"Awesome", CurrentUser:"Zainab"}
}
