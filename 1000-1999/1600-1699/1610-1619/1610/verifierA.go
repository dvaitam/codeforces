package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(cmdPath string, input string) (string, error) {
	cmd := exec.Command(cmdPath)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runRef(input string) (string, error) {
	cmd := exec.Command("go", "run", "1610A.go")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(1)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(1000) + 1
		m := rand.Intn(1000) + 1
		input := fmt.Sprintf("1\n%d %d\n", n, m)
		expected, err := runRef(input)
		if err != nil {
			fmt.Println("failed to run reference:", err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", t, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed: expected %s got %s (input %d %d)\n", t, expected, got, n, m)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
