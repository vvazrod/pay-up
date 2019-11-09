// +build mage

package main

import (
	"fmt"
	"os/exec"
)

func Install() error {
	fmt.Println("Installing dependencies...")
	cmd := exec.Command("go", "mod", "install")
	return cmd.Run()
}

func Test() error {
	fmt.Println("Running tests...")
	cmd := exec.Command("go", "test", "./internal/...", "-coverprofile", "cov.out")
	return cmd.Run()
}

func Coverage() error {
	fmt.Println("Sending coverage data to codecov.io")
	cmd := exec.Command("bash", "<(curl -s https://codecov.io/bash)")
	return cmd.Run()
}
