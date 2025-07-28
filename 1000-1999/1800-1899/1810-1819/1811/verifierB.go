package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

type testB struct {
	n  int64
	x1 int64
	y1 int64
	x2 int64
	y2 int64
}

func genTests() []testB {
	rand.Seed(181102)
	tests := make([]testB, 100)
	for i := range tests {
		n := rand.Int63n(1_000_000) + 2
		if n%2 == 1 {
			n++
		}
		x1 := rand.Int63n(n) + 1
		y1 := rand.Int63n(n) + 1
		x2 := rand.Int63n(n) + 1
		y2 := rand.Int63n(n) + 1
		tests[i] = testB{n, x1, y1, x2, y2}
	}
	return tests
}

func ring(n, x, y int64) int64 {
	if x > n+1-x {
		x = n + 1 - x
	}
	if y > n+1-y {
		y = n + 1 - y
	}
	if x < y {
		return x
	}
	return y
}

func solve(tc testB) int64 {
	r1 := ring(tc.n, tc.x1, tc.y1)
	r2 := ring(tc.n, tc.x2, tc.y2)
	diff := r1 - r2
	if diff < 0 {
		diff = -diff
	}
	return diff
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d %d %d %d\n", tc.n, tc.x1, tc.y1, tc.x2, tc.y2)
	}

	expected := make([]int64, len(tests))
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
		val, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
			os.Exit(1)
		}
		if val != exp {
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
