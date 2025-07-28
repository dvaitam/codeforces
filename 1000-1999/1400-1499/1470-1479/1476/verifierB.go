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
	input  string
	output string
}

func solveCase(p []int64, k int64) int64 {
	sum := p[0]
	var add int64
	for i := 1; i < len(p); i++ {
		required := (p[i]*100 + k - 1) / k
		if sum < required {
			diff := required - sum
			add += diff
			sum += diff
		}
		sum += p[i]
	}
	return add
}

func buildCase(p []int64, k int64) testCase {
	n := len(p)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", p[i]))
	}
	sb.WriteByte('\n')
	output := fmt.Sprintf("%d\n", solveCase(p, k))
	return testCase{input: sb.String(), output: output}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 2
	k := rng.Int63n(10) + 1
	p := make([]int64, n)
	for i := range p {
		p[i] = rng.Int63n(100) + 1
	}
	return buildCase(p, k)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.output)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	cases = append(cases, buildCase([]int64{5, 50, 202, 202}, 1))
	cases = append(cases, buildCase([]int64{1, 1, 1}, 100))
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
