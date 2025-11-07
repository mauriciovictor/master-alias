package utils

import (
	"fmt"
	"os"
)

func Debug(message string, args ...interface{}) {
	if os.Getenv("DEBUG") != "true" {
		return
	}

	fmt.Println(message, args)
}
