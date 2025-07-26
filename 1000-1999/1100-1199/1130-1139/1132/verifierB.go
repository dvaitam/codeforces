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

type testCaseB struct {
	a       []int64
	queries []int
}

func genTestsB() []testCaseB {
	rand.Seed(113202)
	tests := make([]testCaseB, 100)
	for i := range tests {
		n := rand.Intn(5) + 1
		a := make([]int64, n)
		for j := range a {
			a[j] = int64(rand.Intn(20)) + 1
		}
		m := rand.Intn(n) + 1
		q := make([]int, m)
		for j := range q {
			q[j] = rand.Intn(n) + 1
		}
		tests[i] = testCaseB{a: a, queries: q}
	}
	return tests
}

func solveB(tc testCaseB) []int64 {
	n := len(tc.a)
	arr := make([]int64, n)
	copy(arr, tc.a)
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	sum := int64(0)
	for _, v := range arr {
		sum += v
	}
	res := make([]int64, len(tc.queries))
	for i, t := range tc.queries {
		idx := n - t
		res[i] = sum - arr[idx]
	}
	return res
}

func run(bin string, in []byte) ([]byte, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.Bytes(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsB()
	for idx, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintln(&input, len(tc.a))
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		fmt.Fprintln(&input, len(tc.queries))
		for _, q := range tc.queries {
			fmt.Fprintln(&input, q)
		}
		out, err := run(bin, input.Bytes())
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\noutput:\n%s\n", idx+1, err, string(out))
			os.Exit(1)
		}
		scanner := bufio.NewScanner(bytes.NewReader(out))
		scanner.Split(bufio.ScanWords)
		expected := solveB(tc)
		for _, exp := range expected {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", idx+1)
				os.Exit(1)
			}
			val, err := strconv.ParseInt(scanner.Text(), 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", idx+1)
				os.Exit(1)
			}
			if val != exp {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\n", idx+1, exp, val)
				os.Exit(1)
			}
		}
		if scanner.Scan() {
			fmt.Fprintf(os.Stderr, "extra output on test %d\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
