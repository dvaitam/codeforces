package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

var fib [45]int64

func init() {
	fib[0], fib[1] = 1, 1
	for i := 2; i < len(fib); i++ {
		fib[i] = fib[i-1] + fib[i-2]
	}
}

type testD struct {
	n int
	x int64
	y int64
}

func genTests() []testD {
	rand.Seed(181104)
	tests := make([]testD, 100)
	for i := range tests {
		n := rand.Intn(40) + 1
		x := rand.Int63n(fib[n]) + 1
		y := rand.Int63n(fib[n+1]) + 1
		tests[i] = testD{n: n, x: x, y: y}
	}
	return tests
}

func possible(n int, x, y int64) bool {
	for n > 1 {
		if fib[n-1] < y && y <= fib[n] {
			return false
		}
		if y <= fib[n-1] {
			x, y = y, fib[n]-x+1
		} else {
			y -= fib[n]
			x, y = y, fib[n]-x+1
		}
		n--
	}
	return true
}

func solve(tc testD) string {
	if possible(tc.n, tc.x, tc.y) {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d %d\n", tc.n, tc.x, tc.y)
	}

	expected := make([]string, len(tests))
	for i, tc := range tests {
		expected[i] = solve(tc)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, out.String())
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for i, exp := range expected {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
			os.Exit(1)
		}
		got := scanner.Text()
		if got != exp {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
