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

type testCaseD struct {
	input    string
	expected int64
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func solveD(s string, costs []int64) int64 {
	target := "hard"
	inf := int64(1) << 60
	f := make([]int64, 5)
	for i := 1; i <= 4; i++ {
		f[i] = inf
	}
	for i := 0; i < len(s); i++ {
		a := costs[i]
		for j := 3; j >= 0; j-- {
			if s[i] == target[j] {
				if f[j] < f[j+1] {
					f[j+1] = f[j]
				}
				f[j] += a
			}
		}
	}
	res := f[0]
	for j := 1; j <= 3; j++ {
		if f[j] < res {
			res = f[j]
		}
	}
	return res
}

func generateCaseD(rng *rand.Rand) testCaseD {
	n := rng.Intn(10) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	s := string(b)
	costs := make([]int64, n)
	for i := 0; i < n; i++ {
		costs[i] = int64(rng.Intn(10) + 1)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n%s\n", n, s))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", costs[i]))
	}
	sb.WriteByte('\n')
	return testCaseD{input: sb.String(), expected: solveD(s, costs)}
}

func runCaseD(bin string, tc testCaseD) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != tc.expected {
		return fmt.Errorf("expected %d got %d", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	cases := []testCaseD{generateCaseD(rng)}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseD(rng))
	}
	for i, tc := range cases {
		if err := runCaseD(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
