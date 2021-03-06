package main

import "myapp/server"

func main() {
	e := server.NewServer()
	server.ListenAndServe(e)
}
