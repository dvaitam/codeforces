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

type testC struct {
	n int
	b []int64
}

func genTests() []testC {
	rand.Seed(181103)
	tests := make([]testC, 100)
	for i := range tests {
		n := rand.Intn(10) + 2
		a := make([]int64, n)
		for j := range a {
			a[j] = rand.Int63n(1000)
		}
		b := make([]int64, n-1)
		for j := 0; j < n-1; j++ {
			if a[j] > a[j+1] {
				b[j] = a[j]
			} else {
				b[j] = a[j+1]
			}
		}
		tests[i] = testC{n: n, b: b}
	}
	return tests
}

func solve(tc testC) []int64 {
	a := make([]int64, tc.n)
	a[0] = tc.b[0]
	for i := 1; i <= tc.n-2; i++ {
		if tc.b[i] < tc.b[i-1] {
			a[i] = tc.b[i]
		} else {
			a[i] = tc.b[i-1]
		}
	}
	a[tc.n-1] = tc.b[tc.n-2]
	return a
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.n)
		for i, v := range tc.b {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
	}

	expected := make([][]int64, len(tests))
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
	for i, expArr := range expected {
		for j, exp := range expArr {
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
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
