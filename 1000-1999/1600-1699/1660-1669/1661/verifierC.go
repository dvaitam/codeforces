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

func calc(h []int64, H int64) int64 {
	var ones, twos int64
	for _, x := range h {
		diff := H - x
		if diff > 0 {
			ones += diff % 2
			twos += diff / 2
		}
	}
	ans := twos * 2
	if tmp := ones*2 - 1; tmp > ans {
		ans = tmp
	}
	return ans
}

func solveCase(h []int64) int64 {
	mx := h[0]
	for _, v := range h {
		if v > mx {
			mx = v
		}
	}
	d1 := calc(h, mx)
	d2 := calc(h, mx+1)
	if d1 < d2 {
		return d1
	}
	return d2
}

type testCase struct {
	input    string
	expected string
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 100)
	// edge case n=1
	cases = append(cases, func() testCase {
		h := []int64{1}
		var sb strings.Builder
		sb.WriteString("1\n1\n1\n")
		exp := fmt.Sprintf("%d", solveCase(h))
		return testCase{input: sb.String(), expected: exp}
	}())
	for len(cases) < 100 {
		n := rng.Intn(20) + 1
		h := make([]int64, n)
		for i := 0; i < n; i++ {
			h[i] = rng.Int63n(1_000_000_000) + 1
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", h[i]))
		}
		sb.WriteByte('\n')
		exp := fmt.Sprintf("%d", solveCase(h))
		cases = append(cases, testCase{input: sb.String(), expected: exp})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genTests()
	for i, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
