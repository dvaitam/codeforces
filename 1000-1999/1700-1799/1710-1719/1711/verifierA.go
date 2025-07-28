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

type TestCase struct {
	N int
}

func generateTests() []TestCase {
	r := rand.New(rand.NewSource(42))
	tests := make([]TestCase, 100)
	for i := 0; i < 100; i++ {
		n := r.Intn(1000) + 1
		tests[i] = TestCase{N: n}
	}
	return tests
}

func expectedPermutation(n int) []string {
	if n == 1 {
		return []string{"1"}
	}
	res := make([]string, n)
	for i := 2; i <= n; i++ {
		res[i-2] = strconv.Itoa(i)
	}
	res[n-1] = "1"
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, t := range tests {
		fmt.Fprintln(&input, t.N)
	}
	cmd := exec.Command(binary)
	cmd.Stdin = &input
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "binary execution failed:", err)
		os.Exit(1)
	}
	outputs := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(outputs) != len(tests) {
		fmt.Fprintln(os.Stderr, "wrong number of output lines:", len(outputs))
		os.Exit(1)
	}
	for i, line := range outputs {
		fields := strings.Fields(line)
		expect := expectedPermutation(tests[i].N)
		if len(fields) != len(expect) {
			fmt.Fprintf(os.Stderr, "case %d: expected %d numbers, got %d\n", i+1, len(expect), len(fields))
			os.Exit(1)
		}
		for j, val := range expect {
			if fields[j] != val {
				fmt.Fprintf(os.Stderr, "case %d: expected %s at pos %d, got %s\n", i+1, val, j+1, fields[j])
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
