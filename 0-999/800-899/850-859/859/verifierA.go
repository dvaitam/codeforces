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

type CaseA struct {
	input    string
	expected int
}

func expectedA(k int, ranks []int) int {
	maxRank := 0
	for _, r := range ranks {
		if r > maxRank {
			maxRank = r
		}
	}
	if maxRank < 25 {
		return 0
	}
	return maxRank - 25
}

func generateCaseA(rng *rand.Rand) CaseA {
	k := rng.Intn(25) + 1
	m := make(map[int]bool)
	ranks := make([]int, 0, k)
	for len(ranks) < k {
		r := rng.Intn(1_000_000) + 1
		if !m[r] {
			m[r] = true
			ranks = append(ranks, r)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", k))
	for i, r := range ranks {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", r))
	}
	sb.WriteByte('\n')
	return CaseA{sb.String(), expectedA(k, ranks)}
}

func runCase(exe string, input string, expected int) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	var got int
	if _, err := fmt.Sscan(outStr, &got); err != nil {
		return fmt.Errorf("cannot parse output: %v\n%s", err, outStr)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	cases := []CaseA{
		{"25\n1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25\n", 0},
		{"3\n25 26 27\n", 2},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseA(rng))
	}
	for i, tc := range cases {
		if err := runCase(exe, tc.input, tc.expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
