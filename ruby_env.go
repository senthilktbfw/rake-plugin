package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type RubyEnvironmentCollection struct {
	userHomeDir               string
	RubyEnvironmentsMap       map[string]RubyEnvironment
	RubyEnvironmentsGlobalMap map[string]RubyEnvironment
}

func GetNewRubyEnvironmentCollection(userHomeDir string) (RubyEnvironmentCollection, error) {

	rubyEnvCollection := RubyEnvironmentCollection{
		userHomeDir:               userHomeDir,
		RubyEnvironmentsMap:       make(map[string]RubyEnvironment),
		RubyEnvironmentsGlobalMap: make(map[string]RubyEnvironment),
	}

	return rubyEnvCollection, nil
}

func (r *RubyEnvironmentCollection) ToJsonString() (string, error) {
	s, err := StructToJSON(r)
	return s, err
}

func (rec *RubyEnvironmentCollection) GetRvmDir() string {
	rvmDir := filepath.Join(rec.userHomeDir, ".rvm")
	return rvmDir
}

func (r *RubyEnvironmentCollection) GetRubiesDir() string {
	return filepath.Join(r.GetRvmDir(), "rubies")
}

func (r *RubyEnvironmentCollection) GetRubiesForVersion(version string) string {
	return filepath.Join(r.GetRubiesDir(), RubyDirPrefix+version)
}

func (r *RubyEnvironmentCollection) GetGemsDir() string {
	return filepath.Join(r.GetRvmDir(), "gems")
}

func (r *RubyEnvironmentCollection) FindRubyEnvironments() error {

	if !r.IsRmvDirFound() {
		return errors.New("RVM directory not found")
	}

	versionsList, err := r.FindAllInstalledVersions()
	if err != nil {
		return err
	}

	if len(versionsList) == 0 {
		fmt.Println("No Ruby versions found")
		return errors.New("No Ruby versions found")
	}

	for _, version := range versionsList {
		r.AddRubyEnvironmentToCollection(version)
	}

	return nil
}

func (r *RubyEnvironmentCollection) AddRubyEnvironmentToCollection(version string) {

	gemsDir, err := r.IsGemsDirFound(version, false)

	if err == nil {
		rubiesDirForVersion := r.GetRubiesForVersion(version)
		rubyEnv := GetNewRubyEnvironment(version, rubiesDirForVersion, gemsDir, false)
		r.RubyEnvironmentsMap[version] = rubyEnv
	} else {
		fmt.Println(`r.IsGemsDirFound(version, false) failed for version`, version, " err == ", err.Error())
	}

	gemsGlobalDir, err := r.IsGemsDirFound(version, true)
	if err == nil {
		rubiesDirForVersion := r.GetRubiesForVersion(version)
		rubyEnv := GetNewRubyEnvironment(version, rubiesDirForVersion, gemsGlobalDir, true)
		r.RubyEnvironmentsGlobalMap[version] = rubyEnv
	} else {
		fmt.Println(`r.IsGemsDirFound(version, true) failed for version`, version, " err == ", err.Error())
	}

}

func (r *RubyEnvironmentCollection) IsGemsDirFound(version string, isCheckGlobal bool) (string, error) {

	gemsDir := filepath.Join(r.GetGemsDir(), RubyDirPrefix+version)

	if isCheckGlobal {
		gemsDir += "@global"
	}

	_, err := os.Stat(gemsDir)
	if err != nil {
		fmt.Println("Gems directory not found for version", version)
	}

	return gemsDir, err
}

func (r *RubyEnvironmentCollection) FindAllInstalledVersions() ([]string, error) {
	var rubyVersions []string

	rubyPattern := regexp.MustCompile(`^ruby-(\d+\.\d+\.\d+)$`)

	rubiesDir := r.GetRubiesDir()

	rubiesDirList, err := os.ReadDir(rubiesDir)
	if err != nil {
		fmt.Println("FindAllInstalledVersions os.ReadDir failed", err)
		return nil, err
	}

	for _, tmpRubyDir := range rubiesDirList {
		matches := rubyPattern.FindStringSubmatch(tmpRubyDir.Name())
		if len(matches) > 1 {
			rubyVersions = append(rubyVersions, matches[1])
		}
	}

	return rubyVersions, nil
}

func (r *RubyEnvironmentCollection) IsRmvDirFound() bool {
	_, err := os.Stat(r.GetRvmDir())
	return err == nil

}

type RubyEnvironment struct {
	Version     string
	RubyPath    string
	GemPath     string
	IsGlobalGem bool
}

func GetNewRubyEnvironment(version, rubyPath, gemPath string, isGlobalGem bool) RubyEnvironment {
	rubyEnv := RubyEnvironment{
		Version:     version,
		RubyPath:    rubyPath,
		GemPath:     gemPath,
		IsGlobalGem: isGlobalGem,
	}
	return rubyEnv
}
