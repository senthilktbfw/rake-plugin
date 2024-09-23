package main

import (
	"fmt"
	"os"
)

func main() {

	fmt.Println("rake-plugin")

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(BadExit)
	}

	rec, err := GetNewRubyEnvironmentCollection(userHomeDir)
	if err != nil {
		fmt.Println(err)
	}

	err = rec.FindRubyEnvironments()
	if err != nil {
		fmt.Println(err)
		os.Exit(BadExit)
	}

	s, _ := rec.ToJsonString()
	fmt.Println(s):
}
