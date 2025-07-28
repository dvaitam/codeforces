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

type testCaseC struct {
	n int
	k int
	a []int
}

func genTestsC() []testCaseC {
	rng := rand.New(rand.NewSource(44))
	tests := make([]testCaseC, 100)
	for i := range tests {
		n := rng.Intn(8) + 2
		k := rng.Intn(6) + 1
		a := make([]int, n)
		for j := range a {
			a[j] = rng.Intn(10)
		}
		tests[i] = testCaseC{n, k, a}
	}
	return tests
}

func solveC(tc testCaseC) []int {
	exists := make([]bool, tc.k+1)
	pmax := make([]int, tc.n)
	cur := 0
	for i := 0; i < tc.n; i++ {
		v := tc.a[i]
		if v <= tc.k {
			exists[v] = true
		}
		if v > cur {
			cur = v
		}
		pmax[i] = cur
	}
	smax := make([]int, tc.n)
	cur = 0
	for i := tc.n - 1; i >= 0; i-- {
		if tc.a[i] > cur {
			cur = tc.a[i]
		}
		smax[i] = cur
	}
	L := make([]int, tc.k+1)
	idx := 0
	for c := 1; c <= tc.k; c++ {
		for idx < tc.n && pmax[idx] < c {
			idx++
		}
		if idx < tc.n {
			L[c] = idx
		} else {
			L[c] = tc.n
		}
	}
	R := make([]int, tc.k+1)
	idx = tc.n - 1
	for c := 1; c <= tc.k; c++ {
		for idx >= 0 && smax[idx] < c {
			idx--
		}
		if idx >= 0 {
			R[c] = idx
		} else {
			R[c] = -1
		}
	}
	res := make([]int, tc.k)
	for c := 1; c <= tc.k; c++ {
		if !exists[c] {
			res[c-1] = 0
			continue
		}
		res[c-1] = 2 * (R[c] - L[c] + 1)
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsC()
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
	}
	expected := make([][]int, len(tests))
	for i, tc := range tests {
		expected[i] = solveC(tc)
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
	for idx, exp := range expected {
		for j := 0; j < len(exp); j++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", idx+1)
				os.Exit(1)
			}
			val, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", idx+1)
				os.Exit(1)
			}
			if val != exp[j] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d element %d: expected %d got %d\n", idx+1, j+1, exp[j], val)
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
