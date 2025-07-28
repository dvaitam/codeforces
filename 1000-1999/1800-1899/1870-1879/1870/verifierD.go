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

type testCaseD struct {
	n int
	c []int64
	k int64
}

func genTestsD() []testCaseD {
	rng := rand.New(rand.NewSource(45))
	tests := make([]testCaseD, 100)
	for i := range tests {
		n := rng.Intn(8) + 1
		c := make([]int64, n)
		for j := range c {
			c[j] = int64(rng.Intn(20) + 1)
		}
		k := int64(rng.Intn(20) + 1)
		tests[i] = testCaseD{n, c, k}
	}
	return tests
}

func solveD(tc testCaseD) []int64 {
	n := tc.n
	a := make([]int64, n)
	const INF int64 = 1 << 60
	mn := INF
	prefix := int64(0)
	k := tc.k
	for i := n - 1; i >= 0; i-- {
		if tc.c[i] < mn {
			mn = tc.c[i]
		}
		times := k / mn
		prefix += times
		k -= times * mn
		a[i] = prefix
	}
	return a
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
		for i, v := range tc.c {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		fmt.Fprintln(&input, tc.k)
	}
	expected := make([][]int64, len(tests))
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
	for idx, exp := range expected {
		for j := 0; j < len(exp); j++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", idx+1)
				os.Exit(1)
			}
			val, err := strconv.ParseInt(scanner.Text(), 10, 64)
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
