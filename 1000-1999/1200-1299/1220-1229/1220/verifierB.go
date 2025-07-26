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
	n int
	a []int64
	M [][]int64
}

func genTestsB() []testB {
	rand.Seed(122002)
	tests := make([]testB, 100)
	for i := range tests {
		n := rand.Intn(4) + 3 // 3..6
		a := make([]int64, n)
		for j := range a {
			a[j] = int64(rand.Intn(9) + 1)
		}
		M := make([][]int64, n)
		for j := 0; j < n; j++ {
			M[j] = make([]int64, n)
			for k := 0; k < n; k++ {
				if j == k {
					M[j][k] = 0
				} else {
					M[j][k] = a[j] * a[k]
				}
			}
		}
		tests[i] = testB{n: n, a: a, M: M}
	}
	return tests
}

func solveB(tc testB) []int64 {
	return tc.a
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
		fmt.Fprintln(&input, tc.n)
		for j := 0; j < tc.n; j++ {
			for k := 0; k < tc.n; k++ {
				if k > 0 {
					input.WriteByte(' ')
				}
				fmt.Fprint(&input, tc.M[j][k])
			}
			input.WriteByte('\n')
		}
	}

	expected := make([][]int64, len(tests))
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
		for j := 0; j < len(exp); j++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
				os.Exit(1)
			}
			val, err := strconv.ParseInt(scanner.Text(), 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
				os.Exit(1)
			}
			if val != exp[j] {
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
