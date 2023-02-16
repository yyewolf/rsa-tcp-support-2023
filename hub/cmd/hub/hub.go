package main

import (
	"hub"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	s := hub.NewServer()

	err := s.Listen(":8000")
	if err != nil {
		panic(err)
	}
}
