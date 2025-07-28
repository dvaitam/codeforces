package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runBinary(bin string) (string, error) {
	cmd := exec.Command(bin)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierF <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	for t := 1; t <= 100; t++ {
		output, err := runBinary(bin)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n%s\n", t, err, output)
			os.Exit(1)
		}
		if len(output) == 0 {
			fmt.Printf("test %d failed: no output\n", t)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}
