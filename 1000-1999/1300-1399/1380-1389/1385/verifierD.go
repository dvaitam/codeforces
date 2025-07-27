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
	n        int
	s        string
	input    string
	expected string
}

func solveRec(b []byte, l, r int, c byte) int {
	if r-l == 1 {
		if b[l] == c {
			return 0
		}
		return 1
	}
	mid := (l + r) / 2
	cntLeft := 0
	for i := l; i < mid; i++ {
		if b[i] == c {
			cntLeft++
		}
	}
	cntRight := 0
	for i := mid; i < r; i++ {
		if b[i] == c {
			cntRight++
		}
	}
	cost1 := (mid - l - cntLeft) + solveRec(b, mid, r, c+1)
	cost2 := (r - mid - cntRight) + solveRec(b, l, mid, c+1)
	if cost1 < cost2 {
		return cost1
	}
	return cost2
}

func buildCase(s string) testCase {
	n := len(s)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n%s\n", n, s))
	res := solveRec([]byte(s), 0, n, 'a')
	return testCase{n: n, s: s, input: sb.String(), expected: fmt.Sprintf("%d\n", res)}
}

func randomCase(rng *rand.Rand) testCase {
	k := rng.Intn(4) + 1
	n := 1 << k
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return buildCase(string(b))
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
	want := strings.TrimSpace(tc.expected)
	if got != want {
		return fmt.Errorf("expected %q got %q", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		buildCase("a"),
		buildCase("ba"),
	}
	for i := 0; i < 100; i++ {
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
