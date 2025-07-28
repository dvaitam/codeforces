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

type testH struct {
	n   int
	arr []int
}

func genTestsH() []testH {
	rand.Seed(1530008)
	tests := make([]testH, 100)
	for i := range tests {
		n := rand.Intn(20) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(100) + 1
		}
		tests[i] = testH{n: n, arr: arr}
	}
	return tests
}

func lowerBound(a []int, x int) int {
	l, r := 0, len(a)
	for l < r {
		m := (l + r) / 2
		if a[m] < x {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}

func lis(arr []int) int {
	b := make([]int, 0)
	for _, x := range arr {
		i := lowerBound(b, x)
		if i == len(b) {
			b = append(b, x)
		} else {
			b[i] = x
		}
	}
	return len(b)
}

func solveH(tc testH) int {
	best := lis(tc.arr)
	rev := make([]int, len(tc.arr))
	for i := 0; i < len(tc.arr); i++ {
		rev[i] = tc.arr[len(tc.arr)-1-i]
	}
	if v := lis(rev); v > best {
		best = v
	}
	return best
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsH()

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
		expected[i] = solveH(tc)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n%s\n", err, stderr.String())
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
		if err != nil || val != exp {
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
