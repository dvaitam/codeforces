package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runRef(input string) (string, error) {
	cmd := exec.Command("go", "run", "317D.go")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(4)
	for i := 0; i < 100; i++ {
		n := rand.Int63n(1000) + 1
		input := fmt.Sprintf("%d\n", n)
		exp, err := runRef(input)
		if err != nil {
			fmt.Println("reference run error:", err)
			return
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("binary run error on test %d: %v\n", i+1, err)
			return
		}
		if exp != got {
			fmt.Printf("mismatch on test %d\ninput: %sexpected: %s\n got: %s\n", i+1, input, exp, got)
			return
		}
	}
	fmt.Println("all tests passed")
}
