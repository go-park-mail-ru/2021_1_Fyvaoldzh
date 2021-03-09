package main

import "kudago/server"

func main() {
	e := server.NewServer()
	server.ListenAndServe(e)
}
