package main

import (
	"go-demo-mongodb/middleware"
	"log"
)

func main() {
	log.Panic(middleware.Start())
}

