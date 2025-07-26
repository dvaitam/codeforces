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

type testCaseA struct {
	n1, n2, n3, n4 int
}

func genTestsA() []testCaseA {
	rand.Seed(113201)
	tests := make([]testCaseA, 100)
	for i := range tests {
		tests[i] = testCaseA{
			n1: rand.Intn(6),
			n2: rand.Intn(6),
			n3: rand.Intn(6),
			n4: rand.Intn(6),
		}
	}
	return tests
}

func solveA(tc testCaseA) int {
	if tc.n1 != tc.n4 {
		return 0
	}
	if tc.n1 == 0 {
		if tc.n3 != 0 {
			return 0
		}
		return 1
	}
	return 1
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsA()
	for idx, tc := range tests {
		input := fmt.Sprintf("%d %d %d %d\n", tc.n1, tc.n2, tc.n3, tc.n4)
		out, err := run(bin, []byte(input))
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\noutput:\n%s\n", idx+1, err, string(out))
			os.Exit(1)
		}
		scanner := bufio.NewScanner(bytes.NewReader(out))
		scanner.Split(bufio.ScanWords)
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "no output on test %d\n", idx+1)
			os.Exit(1)
		}
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", idx+1)
			os.Exit(1)
		}
		expected := solveA(tc)
		if val != expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\n", idx+1, expected, val)
			os.Exit(1)
		}
		if scanner.Scan() {
			fmt.Fprintf(os.Stderr, "extra output on test %d\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
