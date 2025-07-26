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
	s string
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveD(s string) string {
	res := []byte(s)
	stack := make([]int, 0)
	for i := 0; i < len(s); i++ {
		if s[i] == '1' {
			stack = append(stack, i)
		} else {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		}
	}
	for _, idx := range stack {
		res[idx] = '0'
	}
	return string(res)
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(4))
	tests := make([]testCase, 0, 100)
	fixed := []string{"0", "1", "10", "01", "111000", "101010", "1111", "0000"}
	for _, f := range fixed {
		tests = append(tests, testCase{s: f})
	}
	letters := []byte{'0', '1'}
	for len(tests) < 100 {
		n := rng.Intn(50) + 1
		sb := make([]byte, n)
		for i := 0; i < n; i++ {
			sb[i] = letters[rng.Intn(2)]
		}
		tests = append(tests, testCase{s: string(sb)})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%s\n", t.s)
		expect := solveD(t.s)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expect {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, expect, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
