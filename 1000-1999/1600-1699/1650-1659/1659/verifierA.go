package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin, in string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runRef(in string) (string, error) {
	cmd := exec.Command("go", "run", "1659A.go")
	cmd.Stdin = strings.NewReader(in)
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
		n := rand.Intn(98) + 3
		b := rand.Intn(n/2) + 1
		r := n - b
		input := fmt.Sprintf("1\n%d %d %d\n", n, r, b)
		exp, err := runRef(input)
		if err != nil {
			fmt.Println("reference failed:", err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d exec failed: %v\n", t, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed: expected %s got %s\n", t, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
