package main

import (
	"github.com/go-park-mail-ru/2021_1_Fyvaoldzh/server"
)

func main() {
	e := server.NewServer()
	server.ListenAndServe(e)
}