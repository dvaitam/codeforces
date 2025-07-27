package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveA(n int) string {
	if n%4 == 0 {
		return "YES"
	}
	return "NO"
}

func runBinary(binPath string, input string) (string, error) {
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(1)
	const tests = 100
	for t := 0; t < tests; t++ {
		n := rand.Intn(1000000000-3+1) + 3 // 3..1e9
		input := fmt.Sprintf("1\n%d\n", n)
		expected := solveA(n) + "\n"
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(output) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", t+1, input, expected, output)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
