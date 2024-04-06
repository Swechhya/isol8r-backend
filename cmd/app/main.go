package main

import (
	"github.com/Swechhya/isol8r-backend/router"
	"github.com/Swechhya/isol8r-backend/setup"
)

func main() {
	setup.Setup()

	//setup router
	r := router.SetupRouter()
	r.Run()
}
