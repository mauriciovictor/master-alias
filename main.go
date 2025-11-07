package main

import (
	"github.com/joho/godotenv"
	"master-alias.com/cmd"
)

func main() {
	godotenv.Load()
	cmd.Execute()
}
