package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(42))
	tests := make([]string, 0, 100)
	tests = append(tests, "1 1\n")
	tests = append(tests, "5 2\n")
	tests = append(tests, "10 3\n")
	for len(tests) < 100 {
		n := rng.Intn(50) + 1
		k := rng.Intn(5) + 1
		tests = append(tests, fmt.Sprintf("%d %d\n", n, k))
	}
	return tests
}

func verify(n, k int, out string) error {
	out = strings.TrimSpace(out)
	if len(out) != n {
		return fmt.Errorf("expected length %d, got %d", n, len(out))
	}
	for i := 0; i < n; i++ {
		c := out[i]
		if c < 'a' || c >= byte('a'+k) {
			return fmt.Errorf("invalid character '%c'", c)
		}
	}
	seen := make(map[string]struct{})
	for i := 0; i+1 < n; i++ {
		seen[out[i:i+2]] = struct{}{}
	}
	expected := n - 1
	kk := k * k
	if expected > kk {
		expected = kk
	}
	if len(seen) != expected {
		return fmt.Errorf("expected %d unique pairs, got %d", expected, len(seen))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runCandidate(candidate, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, t)
			os.Exit(1)
		}
		var n, k int
		fmt.Sscanf(t, "%d %d", &n, &k)
		if err := verify(n, k, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, t, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
