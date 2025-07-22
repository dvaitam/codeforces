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

func generateCase(rng *rand.Rand) (string, []string) {
	n := rng.Intn(6) + 2 // 2..7
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	sb := strings.Builder{}
	for i := 0; i < n; i++ {
		sb.WriteRune(letters[rng.Intn(26)])
	}
	s := sb.String()
	m := rng.Intn(5) + 1
	words := make([]string, 0, m)
	seen := map[string]bool{}
	for len(words) < m {
		l := rng.Intn(5) + 1
		sb.Reset()
		for i := 0; i < l; i++ {
			sb.WriteRune(letters[rng.Intn(26)])
		}
		w := sb.String()
		if !seen[w] {
			seen[w] = true
			words = append(words, w)
		}
	}
	return s, words
}

func bruteCount(s string, words []string) int {
	n := len(s)
	set := map[string]bool{}
	for a := 0; a < n; a++ {
		for b := a; b < n; b++ {
			for c := b + 1; c < n; c++ {
				for d := c; d < n; d++ {
					w := s[a:b+1] + s[c:d+1]
					set[w] = true
				}
			}
		}
	}
	count := 0
	for _, w := range words {
		if set[w] {
			count++
		}
	}
	return count
}

func runCase(bin string, s string, words []string, expect int) error {
	var input strings.Builder
	input.WriteString(s + "\n")
	input.WriteString(fmt.Sprintf("%d\n", len(words)))
	for _, w := range words {
		input.WriteString(w + "\n")
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		s, words := generateCase(rng)
		exp := bruteCount(s, words)
		if err := runCase(bin, s, words, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput s=%s words=%v\n", i+1, err, s, words)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
