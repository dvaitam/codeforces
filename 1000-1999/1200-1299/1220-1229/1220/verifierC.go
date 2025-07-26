package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

type testC struct {
	s string
}

func genTestsC() []testC {
	rand.Seed(122003)
	tests := make([]testC, 100)
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	for i := range tests {
		n := rand.Intn(20) + 1
		b := make([]byte, n)
		for j := range b {
			b[j] = letters[rand.Intn(len(letters))]
		}
		tests[i] = testC{s: string(b)}
	}
	return tests
}

func solveC(tc testC) []string {
	res := make([]string, len(tc.s))
	minChar := tc.s[0]
	res[0] = "Mike"
	for i := 1; i < len(tc.s); i++ {
		if tc.s[i] > minChar {
			res[i] = "Ann"
		} else {
			res[i] = "Mike"
		}
		if tc.s[i] < minChar {
			minChar = tc.s[i]
		}
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsC()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.s)
	}

	expected := make([][]string, len(tests))
	for i, tc := range tests {
		expected[i] = solveC(tc)
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
		for j := range exp {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
				os.Exit(1)
			}
			if scanner.Text() != exp[j] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
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
