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

type testA struct {
	s string
}

func genTestsA() []testA {
	rand.Seed(122001)
	tests := make([]testA, 100)
	for i := range tests {
		ones := rand.Intn(5)
		zeros := rand.Intn(5)
		for ones == 0 && zeros == 0 {
			ones = rand.Intn(5)
			zeros = rand.Intn(5)
		}
		b := make([]byte, 0, ones*3+zeros*4)
		for j := 0; j < ones; j++ {
			b = append(b, 'o', 'n', 'e')
		}
		for j := 0; j < zeros; j++ {
			b = append(b, 'z', 'e', 'r', 'o')
		}
		rand.Shuffle(len(b), func(i, j int) { b[i], b[j] = b[j], b[i] })
		tests[i] = testA{s: string(b)}
	}
	return tests
}

func solveA(tc testA) []int {
	ones := 0
	zeros := 0
	for _, ch := range tc.s {
		if ch == 'n' {
			ones++
		} else if ch == 'z' {
			zeros++
		}
	}
	res := make([]int, 0, ones+zeros)
	for i := 0; i < ones; i++ {
		res = append(res, 1)
	}
	for i := 0; i < zeros; i++ {
		res = append(res, 0)
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsA()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d\n%s\n", len(tc.s), tc.s)
	}

	expected := make([][]int, len(tests))
	for i, tc := range tests {
		expected[i] = solveA(tc)
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
			val, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
				os.Exit(1)
			}
			if val != exp[j] {
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
