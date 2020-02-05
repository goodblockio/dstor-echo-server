package main

import (
	"fmt"
	"os"
)

func main() {
	var PORT string
	if os.Getenv("PORT") != "" {
		PORT = os.Getenv("PORT")
		} else {
		PORT = "80"
		}
	fmt.Println("Hello, "+PORT)
}


