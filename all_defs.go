package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const (
	RubyInterpreterName        = "ruby"
	JrubyInterpreterName       = "jruby"
	RubyInstallationGlobalType = "global"
	RubyInstallationLocalType  = "local"
	BadExit                    = 1
	RubyDirPrefix              = "ruby-"
	RakeExeName                = "rake"
	GemsDirName                = "gems"
	RubiesDirName              = "rubies"
	LibRubySuffix              = "lib/ruby/gems"
	GemsSuffix                 = "gems"
	BundleExeName              = "bundle"
)

func StructToJSON(inputStruct interface{}) (string, error) {
	jsonData, err := json.Marshal(inputStruct)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func GetRakeExeName() string {
	return RakeExeName + GetExeFileExtension()
}

func GetBundleExeName() string {
	return BundleExeName + GetExeFileExtension()
}

func GetExeFileExtension() string {
	return ""
}

func GetDefaultRvmPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".rvm"), nil
}

func GetOs() string {
	return runtime.GOOS
}

func IsLinux() bool {
	return GetOs() == "linux"
}

func GetPathSeparator() string {
	if IsLinux() {
		return ":"
	}
	return ";"
}

func GetExecutablePathDirs() ([]string, error) {
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, GetPathSeparator())

	var executablePathsList []string

	for _, p := range paths {
		absPath, err := filepath.Abs(p)
		if err != nil {
			return nil, err
		}
		executablePathsList = append(executablePathsList, absPath)
	}

	return executablePathsList, nil
}

func IsFileOrDirExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func GetRubyVersionDirs(root string) ([]string, error) {

	versionPattern := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	var matchingDirs []string

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() && versionPattern.MatchString(entry.Name()) {
			matchingDirs = append(matchingDirs, entry.Name())
		}
	}

	return matchingDirs, nil
}

func FindRakeGemspecs(dirPath string) ([]string, error) {

	var matchingFiles []string
	versionPattern := regexp.MustCompile(`^rake-\d+\.\d+\.\d+\.gemspec$`)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && versionPattern.MatchString(info.Name()) {
			matchingFiles = append(matchingFiles, path)
		}
		return nil
	})

	return matchingFiles, err
}

//
