package main

import (
	"github.com/bekars/godueros"
)

func main() {
	core, _ := godueros.NewDuCore()
	core.Run()
	return
}
