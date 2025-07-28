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

type TestCase struct {
	n   int
	arr []int
}

func genTests() []TestCase {
	rand.Seed(time.Now().UnixNano())
	const T = 100
	tests := make([]TestCase, T)
	for i := 0; i < T; i++ {
		n := rand.Intn(8) + 2
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(n + 2)
		}
		tests[i] = TestCase{n: n, arr: arr}
	}
	return tests
}

func buildInput(tests []TestCase) []byte {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, len(tests))
	for _, tc := range tests {
		fmt.Fprint(&buf, tc.n)
		for _, v := range tc.arr {
			fmt.Fprint(&buf, " ", v)
		}
		fmt.Fprintln(&buf)
	}
	return buf.Bytes()
}

func runBinary(path string, input []byte) ([]byte, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	return cmd.CombinedOutput()
}

func mexCount(db map[int]int) int {
	m := 0
	for db[m] > 0 {
		m++
	}
	return m
}

func expected(tc TestCase) string {
	n := tc.n
	t := tc.arr
	pref := make([]int, n)
	suf := make([]int, n)
	db := make(map[int]int)
	mex := 0
	for i := 0; i < n; i++ {
		db[t[i]]++
		for db[mex] > 0 {
			mex++
		}
		pref[i] = mex
	}
	for i := range db {
		delete(db, i)
	}
	mex = 0
	for i := n - 1; i >= 0; i-- {
		db[t[i]]++
		for db[mex] > 0 {
			mex++
		}
		suf[i] = mex
	}
	pos := -1
	for i := 0; i+1 < n; i++ {
		if pref[i] == suf[i+1] {
			pos = i
		}
	}
	if pos < 0 {
		return "-1"
	}
	return fmt.Sprintf("2\n1 %d\n%d %d", pos+1, pos+2, n)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	tests := genTests()
	input := buildInput(tests)

	expectedOutputs := make([]string, len(tests))
	for i, tc := range tests {
		expectedOutputs[i] = expected(tc)
	}

	out, err := runBinary(binary, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error executing binary: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	idx := 0
	for i, exp := range expectedOutputs {
		if idx >= len(lines) {
			fmt.Fprintf(os.Stderr, "output too short on test %d\n", i+1)
			os.Exit(1)
		}
		if exp == "-1" {
			if strings.TrimSpace(lines[idx]) != "-1" {
				fmt.Fprintf(os.Stderr, "mismatch on test %d\nexpected -1 got %s\n", i+1, lines[idx])
				os.Exit(1)
			}
			idx++
		} else {
			parts := strings.Split(exp, "\n")
			for _, p := range parts {
				if idx >= len(lines) || strings.TrimSpace(lines[idx]) != strings.TrimSpace(p) {
					fmt.Fprintf(os.Stderr, "mismatch on test %d\nexpected %s got %s\n", i+1, p, lines[idx])
					os.Exit(1)
				}
				idx++
			}
		}
	}
	if idx != len(lines) {
		fmt.Fprintf(os.Stderr, "extra output lines detected")
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
