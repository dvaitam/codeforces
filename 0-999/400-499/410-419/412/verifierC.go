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

// validateC checks candidate output against constraints rather than a fixed expected string.
// Rule per problem: at each position j, consider all non-'?' letters among patterns.
// - If there are letters and all are the same letter L, output must be L.
// - If there are at least two different letters, output must be '?'.
// - If there are no letters (all '?'), output may be any lowercase letter 'a'..'z'.
func validateC(patterns []string, got string) error {
	if len(patterns) == 0 {
		if strings.TrimSpace(got) != "" {
			return fmt.Errorf("expected empty output for empty input")
		}
		return nil
	}
	k := len(patterns[0])
	if len(got) != k {
		return fmt.Errorf("wrong length: expected %d got %d", k, len(got))
	}
	for j := 0; j < k; j++ {
		var letter byte = 0
		conflict := false
		for _, p := range patterns {
			c := p[j]
			if c == '?' {
				continue
			}
			if letter == 0 {
				letter = c
			} else if letter != c {
				conflict = true
				break
			}
		}
		gj := got[j]
		if conflict {
			if gj != '?' {
				return fmt.Errorf("pos %d: expected '?' due to conflict, got %q", j, gj)
			}
		} else if letter == 0 {
			// Any lowercase letter allowed when unconstrained
			if gj < 'a' || gj > 'z' {
				return fmt.Errorf("pos %d: expected lowercase letter for unconstrained position, got %q", j, gj)
			}
		} else {
			if gj != letter {
				return fmt.Errorf("pos %d: expected %q got %q", j, letter, gj)
			}
		}
	}
	return nil
}

func runCase(bin string, patterns []string) error {
	input := fmt.Sprintf("%d\n", len(patterns))
	for _, p := range patterns {
		input += p + "\n"
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if err := validateC(patterns, got); err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

func genCase(rng *rand.Rand) []string {
	n := rng.Intn(6) + 1
	l := rng.Intn(10) + 1
	patterns := make([]string, n)
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	for i := 0; i < n; i++ {
		var sb strings.Builder
		for j := 0; j < l; j++ {
			if rng.Intn(4) == 0 {
				sb.WriteByte('?')
			} else {
				sb.WriteByte(alphabet[rng.Intn(len(alphabet))])
			}
		}
		patterns[i] = sb.String()
	}
	return patterns
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		patterns := genCase(rng)
		if err := runCase(bin, patterns); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%v\n", i+1, err, patterns)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
