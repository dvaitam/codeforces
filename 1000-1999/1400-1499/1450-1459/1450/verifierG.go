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

func solve(n, a, b int, s string) string {
	alpha := 26
	fpos := make([]int, alpha)
	lpos := make([]int, alpha)
	cnt := make([]int, alpha)
	for i := range fpos {
		fpos[i] = n + 1
	}
	for i, ch := range s {
		idx := int(ch - 'a')
		cnt[idx]++
		pos := i + 1
		if pos < fpos[idx] {
			fpos[idx] = pos
		}
		if pos > lpos[idx] {
			lpos[idx] = pos
		}
	}
	good := make([]bool, alpha)
	for i := 0; i < alpha; i++ {
		if cnt[i] == 0 {
			continue
		}
		length := lpos[i] - fpos[i] + 1
		if a*length <= b*cnt[i] {
			good[i] = true
		}
	}
	var ans []rune
	for i := 0; i < alpha; i++ {
		if cnt[i] == 0 {
			continue
		}
		ok := true
		for j := 0; j < alpha; j++ {
			if j == i || cnt[j] == 0 {
				continue
			}
			if !good[j] {
				ok = false
				break
			}
		}
		if ok {
			ans = append(ans, rune('a'+i))
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d", len(ans)))
	if len(ans) > 0 {
		sb.WriteByte(' ')
		for i, ch := range ans {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(string(ch))
		}
	}
	return sb.String()
}

func buildCase(n, a, b int, s string) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n%s\n", n, a, b, s)
	exp := solve(n, a, b, s)
	return testCase{input: sb.String(), expected: exp}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 1
	a := rng.Intn(5) + 1
	b := rng.Intn(5) + 1
	bStr := make([]byte, n)
	for i := 0; i < n; i++ {
		bStr[i] = byte('a' + rng.Intn(26))
	}
	return buildCase(n, a, b, string(bStr))
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
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
