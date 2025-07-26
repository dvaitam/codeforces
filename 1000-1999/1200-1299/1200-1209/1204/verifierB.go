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
	l int
	r int
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

func solveB(n, l, r int) string {
	minSum := int64(n - l)
	minSum += int64(1<<uint(l)) - 1
	maxSum := int64(1<<uint(r)) - 1
	maxSum += int64(n-r) * int64(1<<uint(r-1))
	return fmt.Sprintf("%d %d", minSum, maxSum)
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(2))
	tests := make([]testCase, 0, 100)
	fixed := []testCase{{n: 1, l: 1, r: 1}, {n: 5, l: 2, r: 3}, {n: 10, l: 1, r: 1}, {n: 20, l: 5, r: 5}}
	tests = append(tests, fixed...)
	for len(tests) < 100 {
		n := rng.Intn(1000) + 1
		maxLR := n
		if maxLR > 20 {
			maxLR = 20
		}
		l := rng.Intn(maxLR) + 1
		r := l + rng.Intn(maxLR-l+1)
		if r > maxLR {
			r = maxLR
		}
		tests = append(tests, testCase{n: n, l: l, r: r})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d %d\n", t.n, t.l, t.r)
		expect := solveB(t.n, t.l, t.r)
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
