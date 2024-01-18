package main

import (
	_ "modernc.org/sqlite"
	"todo/model"
	"todo/router"
)

func main() {
	model.Setup()
	router.Setup()
}
