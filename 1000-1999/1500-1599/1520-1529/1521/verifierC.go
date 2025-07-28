package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cmd := exec.Command(bin)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		fmt.Printf("runtime error: %v\n%s\n", err, errb.String())
		os.Exit(1)
	}
	result := strings.TrimSpace(out.String())
	expected := "Problem C is interactive and cannot be automatically solved."
	if result != expected {
		fmt.Printf("expected %q got %q\n", expected, result)
		os.Exit(1)
	}
	fmt.Println("All 1 tests passed")
}
