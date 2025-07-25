package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n      int64
	l      int64
	r      int64
	expect int64
}

var lenMemo map[int64]int64
var onesMemo map[int64]int64

func seqLen(n int64) int64 {
	if n <= 1 {
		return 1
	}
	if v, ok := lenMemo[n]; ok {
		return v
	}
	v := 2*seqLen(n/2) + 1
	lenMemo[n] = v
	return v
}

func onesCount(n int64) int64 {
	if n <= 1 {
		return n
	}
	if v, ok := onesMemo[n]; ok {
		return v
	}
	v := 2*onesCount(n/2) + n%2
	onesMemo[n] = v
	return v
}

func prefix(n, pos int64) int64 {
	if pos <= 0 || n == 0 {
		return 0
	}
	if n == 1 {
		if pos >= 1 {
			return 1
		}
		return 0
	}
	leftLen := seqLen(n / 2)
	if pos <= leftLen {
		return prefix(n/2, pos)
	}
	leftOnes := onesCount(n / 2)
	if pos == leftLen+1 {
		return leftOnes + n%2
	}
	return leftOnes + n%2 + prefix(n/2, pos-leftLen-1)
}

func expectedOnes(n, l, r int64) int64 {
	lenMemo = make(map[int64]int64)
	onesMemo = make(map[int64]int64)
	return prefix(n, r) - prefix(n, l-1)
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 100)
	deterministic := []testCase{
		{n: 0, l: 1, r: 1},
		{n: 1, l: 1, r: 1},
		{n: 7, l: 2, r: 5},
		{n: 10, l: 3, r: 10},
	}
	for _, tc := range deterministic {
		tc.expect = expectedOnes(tc.n, tc.l, tc.r)
		tests = append(tests, tc)
	}
	for len(tests) < 100 {
		n := rng.Int63n(1 << 50) // keep reasonably small
		length := expectedSeqLen(n)
		l := rng.Int63n(length) + 1
		r := l + rng.Int63n(length-l+1)
		tc := testCase{n: n, l: l, r: r}
		tc.expect = expectedOnes(n, l, r)
		tests = append(tests, tc)
	}
	return tests
}

func expectedSeqLen(n int64) int64 {
	lenMemo = make(map[int64]int64)
	return seqLen(n)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("%d %d %d\n", tc.n, tc.l, tc.r)
		expected := strconv.FormatInt(tc.expect, 10)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
