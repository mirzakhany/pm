package main

import "projectmanager/internal/server"

func main() {
	err := server.Start()
	panic(err)
}
