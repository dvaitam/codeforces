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

type testCaseA struct {
	n int64
	k int64
	x int64
}

func genTestsA() []testCaseA {
	rng := rand.New(rand.NewSource(42))
	tests := make([]testCaseA, 100)
	for i := range tests {
		tests[i] = testCaseA{
			n: int64(rng.Intn(200) + 1),
			k: int64(rng.Intn(200) + 1),
			x: int64(rng.Intn(200) + 1),
		}
	}
	return tests
}

func solveA(tc testCaseA) int64 {
	if tc.n < tc.k || tc.x < tc.k-1 {
		return -1
	}
	base := tc.k * (tc.k - 1) / 2
	val := tc.x
	if val == tc.k {
		val--
	}
	return base + (tc.n-tc.k)*val
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsA()
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d %d\n", tc.n, tc.k, tc.x)
	}
	expected := make([]int64, len(tests))
	for i, tc := range tests {
		expected[i] = solveA(tc)
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
