package main

import (

)

func main() {
	e := server.NewServer()
	server.ListenAndServe(e)
}