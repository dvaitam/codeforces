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

type testD struct {
	n   int
	arr []int
}

func genTestsD() []testD {
	rand.Seed(45)
	tests := make([]testD, 100)
	for i := range tests {
		n := rand.Intn(10) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(11) - 5
		}
		tests[i] = testD{n: n, arr: arr}
	}
	return tests
}

func solveD(tc testD) int {
	if tc.n <= 2 {
		return 0
	}
	const inf = int(1e9)
	res := inf
	deltas := []int{-1, 0, 1}
	for _, d1 := range deltas {
		for _, d2 := range deltas {
			start := tc.arr[0] + d1
			diff := (tc.arr[1] + d2) - start
			changes := 0
			if d1 != 0 {
				changes++
			}
			if d2 != 0 {
				changes++
			}
			ok := true
			for i := 2; i < tc.n; i++ {
				expected := start + diff*i
				delta := expected - tc.arr[i]
				if delta < -1 || delta > 1 {
					ok = false
					break
				}
				if delta != 0 {
					changes++
				}
			}
			if ok && changes < res {
				res = changes
			}
		}
	}
	if res == inf {
		return -1
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsD()

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
		expected[i] = solveD(tc)
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
		val, err := strconv.Atoi(scanner.Text())
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
