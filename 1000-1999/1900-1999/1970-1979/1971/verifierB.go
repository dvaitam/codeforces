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

type testCaseB struct{ s string }

func generateCase(rng *rand.Rand) testCaseB {
	n := rng.Intn(10) + 1
	letters := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return testCaseB{string(b)}
}

func checkOutput(tc testCaseB, out string) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("no output")
	}
	ans := strings.ToLower(strings.TrimSpace(lines[0]))
	if ans == "no" {
		// should be impossible
		for i := 1; i < len(tc.s); i++ {
			if tc.s[i] != tc.s[0] {
				return fmt.Errorf("expected YES, got NO")
			}
		}
		return nil
	}
	if ans != "yes" {
		return fmt.Errorf("first line should be YES or NO")
	}
	if len(lines) < 2 {
		return fmt.Errorf("missing permutation line")
	}
	r := strings.TrimSpace(lines[1])
	if len(r) != len(tc.s) {
		return fmt.Errorf("permutation wrong length")
	}
	if r == tc.s {
		return fmt.Errorf("permutation equal to original")
	}
	// check multiset equality
	cnt := make(map[rune]int)
	for _, ch := range tc.s {
		cnt[ch]++
	}
	for _, ch := range r {
		cnt[ch]--
	}
	for _, v := range cnt {
		if v != 0 {
			return fmt.Errorf("permutation mismatch")
		}
	}
	return nil
}

func runCase(bin string, tc testCaseB) error {
	input := fmt.Sprintf("1\n%s\n", tc.s)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return checkOutput(tc, out.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", i+1, err, tc.s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
