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
	input    string
	expected string
}

func solveCase(n int, ranges [][2]int) string {
	cnt := make([]int, n+1)
	for _, r := range ranges {
		for d := r[0]; d <= r[1]; d++ {
			cnt[d]++
		}
	}
	for d := 1; d <= n; d++ {
		if cnt[d] != 1 {
			return fmt.Sprintf("%d %d", d, cnt[d])
		}
	}
	return "OK"
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	m := rng.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	ranges := make([][2]int, m)
	for i := 0; i < m; i++ {
		a := rng.Intn(n) + 1
		b := rng.Intn(n-a+1) + a
		ranges[i] = [2]int{a, b}
		sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
	}
	return testCase{input: sb.String(), expected: solveCase(n, ranges)}
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
	if got != tc.expected {
		return fmt.Errorf("expected %q got %q", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	cases := []testCase{{input: "1 1\n1 1\n", expected: "OK"}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
