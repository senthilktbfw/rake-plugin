package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type RubyEnvironmentCollection struct {
	BaseRvmPath                       string
	UserHomeRubyEnvironmentsMap       map[string]RubyEnvironment
	RubyEnvironmentsWithGlobalGemsMap map[string]RubyEnvironment // global gems ruby-3.3.5@global
	SystemWideRubyInstallationMap     map[string]RubyEnvironment
}

func GetNewRubyEnvironmentCollection(baseRvmPath string) (RubyEnvironmentCollection, error) {

	rubyEnvCollection := RubyEnvironmentCollection{
		BaseRvmPath:                       baseRvmPath,
		UserHomeRubyEnvironmentsMap:       make(map[string]RubyEnvironment),
		RubyEnvironmentsWithGlobalGemsMap: make(map[string]RubyEnvironment),
		SystemWideRubyInstallationMap:     make(map[string]RubyEnvironment),
	}

	return rubyEnvCollection, nil
}

func (r *RubyEnvironmentCollection) ToJsonString() (string, error) {
	s, err := StructToJSON(r)
	return s, err
}

func (r *RubyEnvironmentCollection) GetRvmDir() string {
	return r.BaseRvmPath
}

func (r *RubyEnvironmentCollection) GetRubiesDir() string {
	return filepath.Join(r.GetRvmDir(), RubiesDirName)
}

func (r *RubyEnvironmentCollection) GetRubiesForVersion(version string) string {
	return filepath.Join(r.GetRubiesDir(), RubyDirPrefix+version)
}

func (r *RubyEnvironmentCollection) GetGemsDir() string {
	return filepath.Join(r.GetRvmDir(), GemsDirName)
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

	err = r.AddSystemWideRubyEnvironments()
	if err != nil {
		return err
	}

	return nil
}

func (r *RubyEnvironmentCollection) FindRubyEnvironmentForRake(rakeIdStr string) (RubyEnvironment, error) {
	return RubyEnvironment{}, nil
}

func (r *RubyEnvironmentCollection) AddRubyEnvironmentToCollection(version string) {

	gemsDir, err := r.IsGemsDirFound(version, false)
	if err == nil {
		rubiesDirForVersion := r.GetRubiesForVersion(version)
		rubyEnv := GetNewRubyEnvironment(version, rubiesDirForVersion, gemsDir, false)
		r.UserHomeRubyEnvironmentsMap[version] = rubyEnv
	}

	gemsGlobalDir, err := r.IsGemsDirFound(version, true)
	if err == nil {
		rubiesDirForVersion := r.GetRubiesForVersion(version)
		rubyEnv := GetNewRubyEnvironment(version, rubiesDirForVersion, gemsGlobalDir, true)
		r.RubyEnvironmentsWithGlobalGemsMap[version] = rubyEnv
	}

}

func (r *RubyEnvironmentCollection) AddSystemWideRubyEnvironments() error {

	rubyExeList := []string{RubyInterpreterName, JrubyInterpreterName}

	sysPathDirs, _ := GetExecutablePathDirs()

	for _, rubyExe := range rubyExeList {
		for _, pathDir := range sysPathDirs {
			r.CheckSysPathDirForRubyInstallation(pathDir, rubyExe)
		}
	}

	return nil
}

func (r *RubyEnvironmentCollection) CheckSysPathDirForRubyInstallation(sysPath, rubyExeName string) {

	absoluteRubyPath := filepath.Join(sysPath, rubyExeName)
	_, err := os.Stat(absoluteRubyPath)
	if os.IsNotExist(err) {
		return
	}

	parentDir := filepath.Dir(sysPath)
	possibleLibPath := []string{LibRubySuffix, GemsSuffix}

	for _, libPath := range possibleLibPath {
		rubyLibPath := filepath.Join(parentDir, libPath)

		if !IsFileOrDirExists(rubyLibPath) {
			continue
		}

		versionsList, err := GetRubyVersionDirs(rubyLibPath)
		if err != nil {
			continue
		}

		for _, version := range versionsList {

			if _, found := r.SystemWideRubyInstallationMap[version]; found {
				continue
			}

			specificationsDir := filepath.Join(rubyLibPath, version, "specifications")
			rakeGemSpecs, err := FindRakeGemspecs(specificationsDir)
			if err != nil {
				// fmt.Println("rakeGemSpecs not found")
				continue
			}

			if len(rakeGemSpecs) > 0 {
				gemsPath := filepath.Join(rubyLibPath, version)
				rubyEnv := GetNewRubyEnvironment(version, parentDir, gemsPath, false)
				r.SystemWideRubyInstallationMap[version] = rubyEnv
				return // return on first find
			}
		}

	}
}

func (r *RubyEnvironmentCollection) IsGemsDirFound(version string, isCheckGlobalGem bool) (string, error) {

	gemsDir := filepath.Join(r.GetGemsDir(), RubyDirPrefix+version)

	if isCheckGlobalGem {
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
		// fmt.Println("FindAllInstalledVersions os.ReadDir failed", err)
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

//
//
