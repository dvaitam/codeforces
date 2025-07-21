package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseC struct {
	n int
	k int64
}

func genTestsC() []testCaseC {
	rand.Seed(44)
	tests := make([]testCaseC, 20)
	for i := range tests {
		n := rand.Intn(7) + 2
		var k int64 = int64(rand.Intn(n*n + 1))
		tests[i] = testCaseC{n: n, k: k}
	}
	return tests
}

func possible(n int, k int64) bool {
	if k%2 != 0 {
		return false
	}
	rem := k
	for i := 0; i < n-i-1; i++ {
		j := n - i - 1
		d := int64(j - i)
		if d*2 < rem {
			rem -= d * 2
		} else {
			rem = 0
			break
		}
	}
	return rem == 0
}

func manhattan(p []int) int64 {
	var sum int64
	for i, v := range p {
		sum += int64(math.Abs(float64(v - (i + 1))))
	}
	return sum
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
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
	}

	poss := make([]bool, len(tests))
	for i, tc := range tests {
		poss[i] = possible(tc.n, tc.k)
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
	for i, tc := range tests {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
			os.Exit(1)
		}
		word := strings.ToLower(scanner.Text())
		if word == "no" {
			if poss[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected Yes got No\n", i+1)
				os.Exit(1)
			}
			continue
		}
		if word != "yes" {
			fmt.Fprintf(os.Stderr, "invalid output on test %d: %s\n", i+1, word)
			os.Exit(1)
		}
		if !poss[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected No got Yes\n", i+1)
			os.Exit(1)
		}
		perm := make([]int, tc.n)
		for j := 0; j < tc.n; j++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "not enough numbers for permutation on test %d\n", i+1)
				os.Exit(1)
			}
			val, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid number on test %d\n", i+1)
				os.Exit(1)
			}
			perm[j] = val
		}
		used := make([]bool, tc.n+1)
		for _, v := range perm {
			if v < 1 || v > tc.n || used[v] {
				fmt.Fprintf(os.Stderr, "invalid permutation on test %d\n", i+1)
				os.Exit(1)
			}
			used[v] = true
		}
		if manhattan(perm) != tc.k {
			fmt.Fprintf(os.Stderr, "wrong permutation value on test %d\n", i+1)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
