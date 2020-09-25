package main

import (
	"os"
	"proj/internal/server"
)

func main() {
	err := server.Start()
	if err!=nil{
		panic(err)
	}
	os.Exit(0)
}
