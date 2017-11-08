/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"path/filepath"
)

//Parse CO_CODERPO value, and return lintcode repository URI.
func parseReopEnv(env string) (url string, err error) {
	files := strings.Fields(env)
	if len(files) == 0 {
		return "", fmt.Errorf("CO_CODERPO value is null\n")
	}

	for _, v := range files {
		s := strings.Split(v, "=")
		key, value := s[0], s[1]

		switch key {
		case "coderepo":
			url = value
		default:
			fmt.Fprintf(os.Stdout, "[COUT] Unknown Parameter: [%s]\n", s)
		}
	}
	return url, nil
}

//Git clone the code repository, and process will redirect to system stdout.
func gitClone(repo, dest string) error {
	cmd := exec.Command("git", "clone", repo, dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Git clone error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}

	return nil
}

// get all needed lint python source filepath
func getFilelist(path string)(pyfilelist [] string, err error) {
  reterr := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
    if ( f == nil ) {return nil}
    if f.IsDir() {return nil}
    if strings.Contains(path, ".py") == true{
      pyfilelist = append(pyfilelist, path)
    }
    return nil
    })

  if reterr != nil {
    fmt.Printf("filepath.Walk() returned %v\n", reterr)
		os.Exit(1)
  }
  return pyfilelist, nil
}

//lint source file save in dir:path.Join(os.Getenv("GOPATH"), "src","tmp")
//just show the Path after "GOPATH/src/tmp/"  for user
func userOrigFilePath(lintfilepath string) (originFilePath string){
	start := strings.Index(lintfilepath,"tmp")
	len := len(lintfilepath)
	rstring := []rune(lintfilepath)
	// start is frist "tmp/" start Site,so we plus 4,get path after tmp/
	return string(rstring[start+4:len])
}

func userInputLintConfsave2Rcfile(lintparameter string) error{
	// via bash script save lint config to rcfile(pylint.conf)
	rcfile := os.Getenv("LINTCONFFILE")
	tmpbashfile := path.Join(os.Getenv("GOPATH"), "src","tmp","savelintconf.sh")
	//bashfiel :="./savelintconf.sh"
	//generate-bash script file
	file, err := os.OpenFile(tmpbashfile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
  if err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Save user input parameters err = %s\n", err.Error())
		os.Exit(1)
  }
  defer file.Close()
	byteSlice := []byte("#!/bin/sh\n"+
	                    "pylint " + lintparameter + " --rcfile=./pylint.conf --generate-rcfile > "+ rcfile + "\n")
  _, err = file.Write(byteSlice)
  if err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Savelint parameters generate-bash script file error = %s\n", err.Error())
		os.Exit(1)
  }
  cmd := exec.Command("sh", tmpbashfile)
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  if err := cmd.Run(); err != nil {
    fmt.Fprintf(os.Stderr, "[COUT] Run generate-bash script error: %s\n", err.Error())
    fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
  }
	//del tmp bash script file
	if err := os.Remove(tmpbashfile);err != nil{
		fmt.Fprintf(os.Stdout, "[COUT] Rmove tmp  bash script warrning: %s\n", err.Error())
	}

	return nil
}

//execute pylint
func execPylint(path string) error{
	//get all python source files needed linted
	pyfilelist, err := getFilelist(path)
	if err != nil{
		fmt.Fprintf(os.Stderr, "[COUT] GetFilelist error = %s\n", err.Error())
		os.Exit(1)
	}

	//get pylint rcfile, from environment parameter "CO_CODERPO" in pylintcomponet Dockerfile.
	//user can modify the lint configs in containerops/component/images/python/pylint/src/pylint.conf
	rcfile := os.Getenv("LINTCONFFILE")
	if _, err := os.Stat(rcfile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stdout, "[COUT] lint rcfile not exist, use default config \n")
	}

	//lint every python source file
	for _, table1 := range pyfilelist {
		fmt.Fprintf(os.Stdout, "------------------------------\n")
		fmt.Fprintf(os.Stdout, "[COUT] Start lint file：%s\n", userOrigFilePath(table1))
		cmdpylint := exec.Command("pylint", "--rcfile=" + rcfile, table1)
		cmdpylint.Stdout = os.Stdout
		cmdpylint.Stderr = os.Stderr
		//pylint err return nil when source code hasn't any warrning or err.
		if err := cmdpylint.Run(); err == nil {
			fmt.Fprintf(os.Stderr, "[COUT] %s isn't any warning\n", userOrigFilePath(table1))
		}
		fmt.Fprintf(os.Stdout, "[COUT] End lint file：%s end\n", userOrigFilePath(table1))
	}

	return nil
}

func main() {
	//Get the CO_CODERPO from environment parameter "CO_CODERPO"
	repodata := os.Getenv("CO_CODERPO")
	if len(repodata) == 0 {
		fmt.Fprintf(os.Stderr, "[COUT] %s\n", "The CO_CODERPO value is null.")
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "[COUT] CO_TEST = %s\n", repodata)

	//Parse the CO_CODERPO, get the lintcode repository URI and action
	codeRepo, err := parseReopEnv(repodata)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Parse the CO_CODERPO error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}
	//Create the base path within GOPATH.
	basePath := path.Join(os.Getenv("GOPATH"), "src","tmp")
	//Clone the git repository
	if err := gitClone(codeRepo, basePath); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Clone the code repository error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}

	//Get the CO_CODERPO from environment parameter "CO_CODERPO"
	lintparameterdata := os.Getenv("CO_LINTPARA")
	if len(lintparameterdata) == 0 {
		fmt.Fprintf(os.Stdout, "[COUT] User the default pylint prameter\n")
	}else{
		fmt.Fprintf(os.Stdout, "[COUT] Get input pylint prameter is \"%s\",and other is default\n", lintparameterdata)
		// Save  Save user input lintparameter to rcfile(pylint.conf)
		if err := userInputLintConfsave2Rcfile(lintparameterdata); err !=nil{
			fmt.Fprintf(os.Stderr, "[COUT] Save user input lintparameter error: %s\n", err.Error())
			fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
			os.Exit(1)
		}
	}

	//Execute pylint
	if err := execPylint(basePath); err != nil{
		fmt.Fprintf(os.Stderr, "[COUT] Exec pylint error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}

	//Print result
	fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "true")
	os.Exit(0)
}

// Author: tanhaijun@huawei.com
// 2017-7-17 Create
