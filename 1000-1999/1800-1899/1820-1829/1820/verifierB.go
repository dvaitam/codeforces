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

func solveB(s string) int {
	n := len(s)
	tstr := s + s
	maxRun := 0
	cur := 0
	for i := 0; i < len(tstr); i++ {
		if tstr[i] == '1' {
			cur++
			if cur > maxRun {
				maxRun = cur
			}
		} else {
			cur = 0
		}
	}
	if maxRun == 2*n {
		return n * n
	}
	if maxRun > n {
		maxRun = n
	}
	ans := 0
	for w := 1; w <= maxRun; w++ {
		h := maxRun + 1 - w
		if h > n {
			h = n
		}
		area := w * h
		if area > ans {
			ans = area
		}
	}
	return ans
}

func generateTests() []string {
	rand.Seed(2)
	tests := []string{"0", "1", "01", "10"}
	for len(tests) < 100 {
		n := rand.Intn(50) + 1
		b := make([]byte, n)
		for i := range b {
			if rand.Intn(2) == 0 {
				b[i] = '0'
			} else {
				b[i] = '1'
			}
		}
		tests = append(tests, string(b))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: go run verifierB.go /path/to/binary\n")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, t := range tests {
		fmt.Fprintln(&input, t)
	}

	cmd := exec.Command(binary)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "program failed:", err)
		os.Exit(1)
	}

	outputs := strings.Fields(out.String())
	if len(outputs) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d outputs, got %d\n", len(tests), len(outputs))
		os.Exit(1)
	}

	for i, t := range tests {
		got, err := strconv.Atoi(outputs[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid output on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		want := solveB(t)
		if got != want {
			fmt.Fprintf(os.Stderr, "test %d failed: input %q expected %d got %d\n", i+1, t, want, got)
			os.Exit(1)
		}
	}

	fmt.Println("all tests passed")
}
