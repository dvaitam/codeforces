package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseC struct {
	n   int
	k   int64
	arr []int
}

func generateTestsC(num int) []testCaseC {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testCaseC, num)
	for i := 0; i < num; i++ {
		n := rand.Intn(20) + 1
		k := rand.Int63n(1000)
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(n + 1)
		}
		tests[i] = testCaseC{n: n, k: k, arr: arr}
	}
	return tests
}

func solveC(tc testCaseC) []int {
	n := tc.n
	k := tc.k
	a := make([]int, n)
	copy(a, tc.arr)
	r := int(k % int64(n+1))
	freq := make([]int, n+1)
	for _, v := range a {
		freq[v] = 1
	}
	mex := 0
	for mex <= n && freq[mex] == 1 {
		mex++
	}
	x := make([]int, r+1)
	for step := 1; step <= r; step++ {
		x[step] = mex
		leave := a[n-step]
		freq[leave] = 0
		freq[mex] = 1
		if leave < mex {
			mex = leave
		} else {
			for mex <= n && freq[mex] == 1 {
				mex++
			}
		}
	}
	res := make([]int, n)
	for i := 0; i < r; i++ {
		res[i] = x[r-i]
	}
	for i := r; i < n; i++ {
		res[i] = a[i-r]
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsC(100)
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
		for i, v := range tc.arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("binary execution failed:", err)
		os.Exit(1)
	}
	outputs := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(outputs) != len(tests) {
		fmt.Printf("expected %d lines of output, got %d\n", len(tests), len(outputs))
		os.Exit(1)
	}
	for i, tc := range tests {
		expectedArr := solveC(tc)
		expected := strings.TrimSpace(strings.TrimSpace(fmt.Sprint(expectedArr)))
		actual := strings.TrimSpace(outputs[i])
		expected = strings.Join(strings.Fields(expected), " ")
		if actual != expected {
			fmt.Printf("mismatch on test %d:\nexpected: %s\nactual:   %s\n", i+1, expected, actual)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
