package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

type testB struct {
	n   int
	arr []int64
}

func solve(arr []int64) int64 {
	n := int64(len(arr))
	var sum int64
	for _, v := range arr {
		sum += v
	}
	r := sum % n
	return r * (n - r)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(42)
	const T = 100
	tests := make([]testB, T)
	for i := 0; i < T; i++ {
		n := rand.Intn(50) + 1
		arr := make([]int64, n)
		for j := range arr {
			arr[j] = rand.Int63n(1e9)
		}
		tests[i] = testB{n: n, arr: arr}
	}
	var input bytes.Buffer
	fmt.Fprintln(&input, T)
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.n)
		for j, v := range tc.arr {
			if j+1 == len(tc.arr) {
				fmt.Fprintf(&input, "%d\n", v)
			} else {
				fmt.Fprintf(&input, "%d ", v)
			}
		}
	}
	expected := make([]int64, T)
	for i, tc := range tests {
		expected[i] = solve(tc.arr)
	}
	cmd := exec.Command(binary)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for i := 0; i < T; i++ {
		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, "insufficient output")
			os.Exit(1)
		}
		gotStr := scanner.Text()
		var got int64
		fmt.Sscan(gotStr, &got)
		if got != expected[i] {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d got %d\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output after", T, "tests")
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
