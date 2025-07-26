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

type testB struct {
	n int
	s string
}

func genTestsB() []testB {
	rand.Seed(43)
	tests := make([]testB, 100)
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	for i := range tests {
		n := rand.Intn(20) + 1
		b := make([]byte, n)
		for j := range b {
			if rand.Intn(4) == 0 {
				b[j] = 'x'
			} else {
				b[j] = letters[rand.Intn(len(letters))]
			}
		}
		tests[i] = testB{n: n, s: string(b)}
	}
	return tests
}

func solveB(tc testB) int {
	cnt := 0
	run := 0
	for i := 0; i < tc.n; i++ {
		if tc.s[i] == 'x' {
			run++
			if run >= 3 {
				cnt++
			}
		} else {
			run = 0
		}
	}
	return cnt
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
		fmt.Fprintf(&input, "%d\n%s\n", tc.n, tc.s)
	}

	expected := make([]int, len(tests))
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
