package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testA struct {
	n int64
	m int64
}

func genTestsA() []testA {
	rand.Seed(1353)
	tests := make([]testA, 0, 100)
	// edge cases
	tests = append(tests, testA{1, 0}, testA{2, 5}, testA{3, 7}, testA{2, 0})
	tests = append(tests, testA{1, 1000000000}, testA{2, 1000000000}, testA{3, 1000000000})
	for len(tests) < 100 {
		n := rand.Int63n(1000) + 1
		m := rand.Int63n(1000000000)
		tests = append(tests, testA{n, m})
	}
	return tests[:100]
}

func expectedA(n, m int64) int64 {
	if n == 1 {
		return 0
	}
	if n == 2 {
		return m
	}
	return 2 * m
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTestsA()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d %d\n", tc.n, tc.m)
		exp := fmt.Sprintf("%d", expectedA(tc.n, tc.m))
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
