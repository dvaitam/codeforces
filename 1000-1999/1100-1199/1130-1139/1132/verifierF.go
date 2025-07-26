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

type testCaseF struct {
	s string
}

func genTestsF() []testCaseF {
	rand.Seed(113206)
	tests := make([]testCaseF, 100)
	letters := []byte{'a', 'b', 'c'}
	for i := range tests {
		n := rand.Intn(8) + 1
		b := make([]byte, n)
		for j := range b {
			b[j] = letters[rand.Intn(len(letters))]
		}
		tests[i] = testCaseF{s: string(b)}
	}
	return tests
}

func solveF(tc testCaseF) int {
	s := []byte(tc.s)
	seq := []byte{}
	for i := 0; i < len(s); i++ {
		if len(seq) == 0 || seq[len(seq)-1] != s[i] {
			seq = append(seq, s[i])
		}
	}
	s = seq
	n := len(s)
	if n == 0 {
		return 0
	}
	const inf = int(1e9)
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, n)
		for j := range dp[i] {
			dp[i][j] = inf
		}
		dp[i][i] = 1
	}
	for i := n - 1; i >= 0; i-- {
		for j := i + 1; j < n; j++ {
			dp[i][j] = dp[i+1][j] + 1
			for k := i + 1; k <= j; k++ {
				if s[i] == s[k] {
					cost := dp[k][j]
					if k > i+1 {
						cost += dp[i+1][k-1]
					}
					if cost < dp[i][j] {
						dp[i][j] = cost
					}
				}
			}
		}
	}
	return dp[0][n-1]
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsF()
	for idx, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintln(&input, len(tc.s))
		fmt.Fprintln(&input, tc.s)
		out, err := run(bin, input.Bytes())
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
		expected := solveF(tc)
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
