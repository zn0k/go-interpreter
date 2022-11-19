package main

import (
	"fmt"
	"interpreter/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello, %s! Interpreter REPL\n\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
