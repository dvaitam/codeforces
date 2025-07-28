package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

type testA struct {
	n int
	d int
	s string
}

func genTests() []testA {
	rand.Seed(181101)
	tests := make([]testA, 100)
	for i := range tests {
		n := rand.Intn(20) + 1
		d := rand.Intn(10)
		digits := make([]byte, n)
		digits[0] = byte('1' + rand.Intn(9))
		for j := 1; j < n; j++ {
			digits[j] = byte('0' + rand.Intn(10))
		}
		tests[i] = testA{n: n, d: d, s: string(digits)}
	}
	return tests
}

func solve(tc testA) string {
	res := make([]byte, 0, tc.n+1)
	inserted := false
	for i := 0; i < tc.n; i++ {
		if !inserted && int(tc.s[i]-'0') < tc.d {
			res = append(res, byte('0'+tc.d))
			inserted = true
		}
		res = append(res, tc.s[i])
	}
	if !inserted {
		res = append(res, byte('0'+tc.d))
	}
	return string(res)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n%s\n", tc.n, tc.d, tc.s)
	}

	expected := make([]string, len(tests))
	for i, tc := range tests {
		expected[i] = solve(tc)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Env = append(os.Environ(), "GOMAXPROCS=1")
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
		got := scanner.Text()
		if got != exp {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %s got %s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
