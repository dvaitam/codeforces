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

func solveCase(n, k int, b []int64) int64 {
	diff1 := make([]int64, n+2)
	diff2 := make([]int64, n+2)
	var cur1, cur2 int64
	var ans int64
	for i := n; i >= 1; i-- {
		cur1 += diff1[i]
		cur2 += diff2[i]
		val := cur1*int64(i) + cur2
		if val < b[i-1] {
			delta := b[i-1] - val
			if i >= k {
				x := (delta + int64(k) - 1) / int64(k)
				ans += x
				l := i - k + 1
				cur1 += x
				cur2 += x * (1 - int64(l))
				diff1[l-1] -= x
				diff2[l-1] -= x * (1 - int64(l))
			} else {
				x := (delta + int64(i) - 1) / int64(i)
				ans += x
				cur1 += x
			}
		}
	}
	return ans
}

type testCase struct {
	input    string
	expected string
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 100)
	// simple case n=1
	cases = append(cases, func() testCase {
		n := 1
		k := 1
		b := []int64{1}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n%d\n", n, k, b[0]))
		exp := fmt.Sprintf("%d", solveCase(n, k, b))
		return testCase{input: sb.String(), expected: exp}
	}())
	for len(cases) < 100 {
		n := rng.Intn(20) + 1
		k := rng.Intn(n) + 1
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			b[i] = rng.Int63n(1_000_000_000) + 1
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", b[i]))
		}
		sb.WriteByte('\n')
		exp := fmt.Sprintf("%d", solveCase(n, k, b))
		cases = append(cases, testCase{input: sb.String(), expected: exp})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
