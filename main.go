package main

import (
	"fmt"

	"github.com/joshalling/gatorcli/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v", err)
	}

	err = c.SetUser("joshalling")
	if err != nil {
		fmt.Printf("Error writing to config: %v", err)
	}

	c, err = config.Read()
	if err != nil {
		fmt.Printf("Error reading config after write: %v", err)
	}
	fmt.Println(c)
}
