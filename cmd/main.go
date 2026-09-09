package main

import (
	"fmt"
	"os"

	"Containers/internal/container"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage giutainer run <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		container.Run()
	case "child":
		container.Child()
	default:
		fmt.Println("Unknown command. Usage: giutainer run <command>")
		os.Exit(1)
	}
}
