// Copyright 2019 Jeff Nickoloff "jeff@allingeek.com"
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// imitations under the License.
package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"syscall"
)

func main() {
	pairs, command, err := ValidateAndCapture()
	if err != nil {
		os.Stderr.Write([]byte(err.Error() + "\n"))
		os.Exit(1)
		return
	}
	for _, p := range pairs {
		err = Pipe(p.EnvVarName, p.FileName)
		if err != nil {
			os.Stderr.Write([]byte(err.Error() + "\n"))
			os.Exit(2)
			return
		}
	}
	for _, p := range pairs {
		os.Unsetenv(p.EnvVarName)
	}
	Exec(command[0], command[0:])
}

type SecretPair struct {
	EnvVarName string
	FileName   string
}

func ValidateAndCapture() ([]SecretPair, []string, error) {
	pairs := []SecretPair{}
	re := regexp.MustCompile(`(.+?)=(.+)`)
	var rest []string
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == `--` {
			if i >= len(os.Args) {
				return []SecretPair{}, []string{}, errors.New(`no command specified`)
			}
			rest = os.Args[i+1:]
			break
		}
		parts := re.FindSubmatch([]byte(os.Args[i]))
		if len(parts) != 3 {
			return []SecretPair{}, []string{}, errors.New(`invalid construction: ` + os.Args[i])
		}
		pairs = append(pairs, SecretPair{EnvVarName: string(parts[1]), FileName: string(parts[2])})
	}
	return pairs, rest, nil
}

func Pipe(i, o string) error {
	os.MkdirAll(filepath.Dir(o), 0777)

	f, err := os.OpenFile(o, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	envVarValue := os.Getenv(i)
	if len(envVarValue) <= 0 {
		return errors.New("Input variable: " + i + " is unset: " + envVarValue)
	}

	n, err := f.Write([]byte(envVarValue))
	if err == nil && n < len(envVarValue) {
		return errors.New("unable to write complete value to file")
	}
	if err != nil {
		return err
	}
	return nil
}

func Exec(command string, rest []string) {
	name, err := exec.LookPath(command)
	if err != nil {
		os.Stderr.Write([]byte("error: no such executable: " + err.Error() + "\n"))
		os.Exit(5)
		return
	}

	if err = syscall.Exec(name, rest, os.Environ()); err != nil {
		os.Stderr.Write([]byte("error: exec failed: " + err.Error() + "\n"))
		os.Exit(6)
		return
	}
}
