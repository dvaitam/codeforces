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
	n   int
	arr []int
}

func genTestsA() []testCaseA {
	rand.Seed(42)
	tests := make([]testCaseA, 20)
	for i := range tests {
		n := rand.Intn(9) + 2 // 2..10
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(100) + 1
		}
		tests[i] = testCaseA{n, arr}
	}
	return tests
}

func solveA(tc testCaseA) int {
	maxPrev := 0
	for i := 0; i < tc.n-1; i++ {
		if tc.arr[i] > maxPrev {
			maxPrev = tc.arr[i]
		}
	}
	return tc.arr[tc.n-1] + maxPrev
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
		fmt.Fprintln(&input, tc.n)
		for i, v := range tc.arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
	}

	expected := make([]int, len(tests))
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
		val, err := strconv.Atoi(valStr)
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
