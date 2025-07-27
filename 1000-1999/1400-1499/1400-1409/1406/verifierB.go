package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
)

type testB struct {
	n   int
	arr []int64
}

func genTestsB() []testB {
	rand.Seed(2)
	tests := make([]testB, 100)
	for i := range tests {
		n := rand.Intn(10) + 5
		arr := make([]int64, n)
		for j := range arr {
			arr[j] = int64(rand.Intn(6001) - 3000)
		}
		tests[i] = testB{n: n, arr: arr}
	}
	return tests
}

func solveB(tc testB) int64 {
	a := append([]int64(nil), tc.arr...)
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	n := len(a)
	p1 := a[n-1] * a[n-2] * a[n-3] * a[n-4] * a[n-5]
	p2 := a[0] * a[1] * a[n-1] * a[n-2] * a[n-3]
	p3 := a[0] * a[1] * a[2] * a[3] * a[n-1]
	res := p1
	if p2 > res {
		res = p2
	}
	if p3 > res {
		res = p3
	}
	return res
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
		for i, v := range tc.arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
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
