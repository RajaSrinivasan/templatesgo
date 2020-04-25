package cmd

import (
	"fmt"
	"testing"
)

func TestCopyDir(t *testing.T) {
	copyDir("../../templates", "/tmp")
}

func TestCopyFile(t *testing.T) {
	fmt.Println("TestcopyFile")
	copyFile("./", "abc", "/tmp")
}
