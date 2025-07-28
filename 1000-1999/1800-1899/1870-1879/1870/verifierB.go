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
	n int
	m int
	a []int
	b []int
}

func genTestsB() []testCaseB {
	rng := rand.New(rand.NewSource(43))
	tests := make([]testCaseB, 100)
	for i := range tests {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		a := make([]int, n)
		for j := range a {
			a[j] = rng.Intn(50)
		}
		b := make([]int, m)
		for j := range b {
			b[j] = rng.Intn(50)
		}
		tests[i] = testCaseB{n, m, a, b}
	}
	return tests
}

func solveB(tc testCaseB) (int, int) {
	orB := 0
	for _, v := range tc.b {
		orB |= v
	}
	xorA := 0
	for _, v := range tc.a {
		xorA ^= v
	}
	xorOR := 0
	for _, v := range tc.a {
		xorOR ^= (v | orB)
	}
	if tc.n%2 == 1 {
		return xorA, xorOR
	}
	return xorOR, xorA
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
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
	}
	expected := make([][2]int, len(tests))
	for i, tc := range tests {
		x, y := solveB(tc)
		expected[i] = [2]int{x, y}
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
		vals := make([]int, 2)
		for j := 0; j < 2; j++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
				os.Exit(1)
			}
			v, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
				os.Exit(1)
			}
			vals[j] = v
		}
		if vals[0] != exp[0] || vals[1] != exp[1] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d %d got %d %d\n", i+1, exp[0], exp[1], vals[0], vals[1])
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
