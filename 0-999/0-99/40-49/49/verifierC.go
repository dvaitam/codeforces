package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	n int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d\n", t.n)
		expect := solveC(t.n)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, strings.TrimSpace(expect), strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func solveC(n int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d ", n))
	for i := 1; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d ", i))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateTests() []testCase {
	rand.Seed(3)
	tests := make([]testCase, 0, 100)
	fixed := []int{1, 2, 3, 4, 5, 10, 50, 100}
	for _, f := range fixed {
		tests = append(tests, testCase{f})
	}
	for len(tests) < 100 {
		n := rand.Intn(100) + 1
		tests = append(tests, testCase{n})
	}
	return tests
}
