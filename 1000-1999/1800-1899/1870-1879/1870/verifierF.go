package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
)

type testCaseF struct {
	n int
	k int
}

func genTestsF() []testCaseF {
	rng := rand.New(rand.NewSource(47))
	tests := make([]testCaseF, 100)
	for i := range tests {
		n := rng.Intn(7) + 1
		k := rng.Intn(5) + 2
		tests[i] = testCaseF{n, k}
	}
	return tests
}

func baseK(x int, k int) string {
	if x == 0 {
		return "0"
	}
	digits := []byte{}
	for x > 0 {
		digits = append(digits, byte('0'+x%k))
		x /= k
	}
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	return string(digits)
}

func solveF(tc testCaseF) int {
	repr := make([]struct {
		str string
		idx int
	}, tc.n)
	for i := 1; i <= tc.n; i++ {
		repr[i-1] = struct {
			str string
			idx int
		}{baseK(i, tc.k), i}
	}
	sort.Slice(repr, func(i, j int) bool { return repr[i].str < repr[j].str })
	cnt := 0
	for pos, v := range repr {
		if v.idx == pos+1 {
			cnt++
		}
	}
	return cnt
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsF()
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
	}
	expected := make([]int, len(tests))
	for i, tc := range tests {
		expected[i] = solveF(tc)
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
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\n", i+1, exp, val)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
