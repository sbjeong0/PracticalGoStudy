package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"testing"
)

func TestMainApplication(t *testing.T) {
	process := exec.Command("go", "run", "main.go", "1")
	stdin, err := process.StdinPipe()
	if err != nil {
		fmt.Println("failed")
		fmt.Println(err)
		return
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "subin\n")
	}()

	out, err := process.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	if string(out) != "Your name please? Press the Enter key when done\nNice to meet you subin\n" {
		t.Error("Expected output not found")
	}
}
