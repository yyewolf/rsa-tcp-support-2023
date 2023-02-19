package main

import (
	"bot"

	"github.com/joho/godotenv"
)

func main() {
	// Load env
	godotenv.Load()
	// Start agent
	err := bot.Start()
	if err != nil {
		panic(err)
	}
}
