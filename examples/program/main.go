package main

import (
	"fmt"

	"github.com/switchupcb/copygen/cli"
)

// main is run from /copygen.
func main() {
	env := cli.Environment{
		YMLPath: "examples/main/setup/setup.yml",
		Output:  false, // Don't output to standard output.
		Write:   true,  // Write the output to a file.
	}

	code, err := env.Run()
	if err != nil {
		fmt.Printf("%v", err)

		return
	}

	// Print the code to standard output anyways.
	fmt.Println(code)
}
