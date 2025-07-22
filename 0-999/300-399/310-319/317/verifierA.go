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
	cmd := exec.Command("go", "run", "317A.go")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	for i := 0; i < 100; i++ {
		x := rand.Int63n(200) - 100
		y := rand.Int63n(200) - 100
		m := rand.Int63n(200) - 100
		input := fmt.Sprintf("%d %d %d\n", x, y, m)
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
