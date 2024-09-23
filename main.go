package main

import (
	"fmt"
	"os"
)

func main() {

	fmt.Println("rake-plugin")

	defaultRvmPath, err := GetDefaultRvmPath()
	if err != nil {
		fmt.Println(err)
		os.Exit(BadExit)
	}

	rec, err := GetNewRubyEnvironmentCollection(defaultRvmPath)
	if err != nil {
		fmt.Println(err)
	}

	err = rec.FindRubyEnvironments()
	if err != nil {
		fmt.Println(err)
		os.Exit(BadExit)
	}

	rubyEnv, err := rec.FindRubyEnvironmentForRake("2.7.2")
	rakeExecContext, err := GetRakeExecParamsFromArgs(rubyEnv)
	if err != nil {
		fmt.Println(err)
		os.Exit(BadExit)
	}

	err = rakeExecContext.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(BadExit)
	}

	s, _ := rec.ToJsonString()
	fmt.Println(s)
}

func GetRakeExecParamsFromArgs(rubyEnv RubyEnvironment) (RakeExecContext, error) {
	argsList := os.Args
	_ = argsList

	return RakeExecContext{}, nil
}

type RakeExecContext struct {
	RubyEnvInfo      RubyEnvironment
	RakeInstallation string
	RakeFile         string
	RakeLibDir       string
	RakeWorkingDir   string
	Tasks            string
	IsSilent         bool
	IsBundleExec     bool
}

func (r *RakeExecContext) Run() error {

	//exeName := GetRakeExeName()
	//if r.IsBundleExec {
	//	exeName = GetBundleExeName()
	//}
	//cmd := exec.Command(exeName,
	//	"--rakefile", r.RakeFile,
	//	"--libdir", r.RakeLibDir,
	//	"--dir", r.RakeWorkingDir,
	//	r.Tasks)
	//err := cmd.Run()
	//if err != nil {
	//	return err
	//}
	return nil
}

//
//
