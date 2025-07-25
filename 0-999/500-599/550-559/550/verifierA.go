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

func containsNonOverlap(s, p1, p2 string) bool {
	n := len(s)
	for i := 0; i+len(p1) <= n; i++ {
		if s[i:i+len(p1)] == p1 {
			for j := i + len(p1); j+len(p2) <= n; j++ {
				if s[j:j+len(p2)] == p2 {
					return true
				}
			}
			break
		}
	}
	return false
}

func expected(s string) string {
	if containsNonOverlap(s, "AB", "BA") || containsNonOverlap(s, "BA", "AB") {
		return "YES"
	}
	return "NO"
}

func randString(rng *rand.Rand, n int) string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	s := randString(rng, n)
	input := s + "\n"
	exp := expected(s)
	return input, exp
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
