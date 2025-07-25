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

type testCase struct {
	n      int
	piles  []int
	expect string
}

func expectedResult(piles []int) string {
	grundy := make([]int, 61)
	for s := 1; s <= 60; s++ {
		g := 1
		for (g+1)*(g+2)/2 <= s {
			g++
		}
		grundy[s] = g
	}
	xor := 0
	for _, x := range piles {
		xor ^= grundy[x]
	}
	if xor == 0 {
		return "YES"
	}
	return "NO"
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 100)
	deterministic := []testCase{
		{n: 1, piles: []int{1}},
		{n: 2, piles: []int{1, 2}},
		{n: 3, piles: []int{3, 3, 3}},
	}
	for _, tc := range deterministic {
		tc.expect = expectedResult(tc.piles)
		tests = append(tests, tc)
	}
	for len(tests) < 100 {
		n := rng.Intn(10) + 1
		piles := make([]int, n)
		for i := 0; i < n; i++ {
			piles[i] = rng.Intn(60) + 1
		}
		tests = append(tests, testCase{n: n, piles: piles, expect: expectedResult(piles)})
	}
	return tests
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", tc.n))
		for _, v := range tc.piles {
			input.WriteString(fmt.Sprintf("%d ", v))
		}
		input.WriteByte('\n')
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expect, got, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
