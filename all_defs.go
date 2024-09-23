package main

import "encoding/json"

const (
	RubyInterpreterName        = "ruby"
	RubyInstallationGlobalType = "global"
	RubyInstallationLocalType  = "local"
	BadExit                    = 1
	RubyDirPrefix              = "ruby-"
)

func StructToJSON(inputStruct interface{}) (string, error) {
	// Marshal the struct into JSON
	jsonData, err := json.Marshal(inputStruct)
	if err != nil {
		return "", err
	}

	// Convert the JSON byte slice to a string
	return string(jsonData), nil
}

//
