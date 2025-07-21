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

type testCaseB struct {
	n int64
	a int64
	b int64
}

func genTestsB() []testCaseB {
	rand.Seed(43)
	tests := make([]testCaseB, 20)
	for i := range tests {
		tests[i] = testCaseB{
			n: int64(rand.Intn(50) + 1),
			a: int64(rand.Intn(50) + 1),
			b: int64(rand.Intn(50) + 1),
		}
	}
	return tests
}

func profit(n, a, b, k int64) int64 {
	return k*b - k*(k-1)/2 + (n-k)*a
}

func solveB(tc testCaseB) int64 {
	maxk := tc.n
	if tc.b < maxk {
		maxk = tc.b
	}
	if tc.b <= tc.a {
		return tc.n * tc.a
	}
	k := tc.b - tc.a
	if k > maxk {
		k = maxk
	}
	best := profit(tc.n, tc.a, tc.b, k)
	last := profit(tc.n, tc.a, tc.b, maxk)
	if last > best {
		best = last
	}
	return best
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsB()
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d %d\n", tc.n, tc.a, tc.b)
	}
	expected := make([]int64, len(tests))
	for i, tc := range tests {
		expected[i] = solveB(tc)
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
		valStr := scanner.Text()
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "non-integer output on test %d: %s\n", i+1, valStr)
			os.Exit(1)
		}
		if val != exp {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\n", i+1, exp, val)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
