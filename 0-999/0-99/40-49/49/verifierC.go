package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
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
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if !validateOutput(t.n, out) {
			fmt.Printf("test %d failed: invalid permutation for n=%d\n", i+1, t.n)
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

func validateOutput(n int, out string) bool {
	fields := strings.Fields(out)
	if len(fields) != n {
		return false
	}
	p := make([]int, n)
	seen := make([]bool, n+1)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil || v < 1 || v > n || seen[v] {
			return false
		}
		seen[v] = true
		p[i] = v
	}
	for i := 2; i <= n; i++ {
		for j := i; j <= n; j += i {
			if p[j-1]%i == 0 {
				return false
			}
		}
	}
	return true
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
