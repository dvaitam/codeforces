package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	oracle := filepath.Join(os.TempDir(), fmt.Sprintf("oracle58D_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", oracle, refSrc)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	// n must be even, 2 <= n <= 10 for small tests
	n := (rng.Intn(4) + 1) * 2 // 2, 4, 6, or 8

	// We need names such that a valid calendar exists.
	// All lines must have the same length = len(name1) + 1 + len(name2).
	// So the sum of lengths of each pair must be constant.
	// Simplest: generate names all of the same length, or two groups with complementary lengths.
	// Use two lengths that sum to a constant.
	minLen := rng.Intn(3) + 1
	maxLen := rng.Intn(3) + minLen + 1

	// Generate n/2 names of minLen and n/2 names of maxLen
	nameSet := make(map[string]bool)
	names := make([]string, 0, n)
	for len(names) < n/2 {
		s := randName(rng, minLen)
		if !nameSet[s] {
			nameSet[s] = true
			names = append(names, s)
		}
	}
	for len(names) < n {
		s := randName(rng, maxLen)
		if !nameSet[s] {
			nameSet[s] = true
			names = append(names, s)
		}
	}

	// Shuffle names
	rng.Shuffle(len(names), func(i, j int) { names[i], names[j] = names[j], names[i] })

	// Pick separator: ASCII 33-126, excluding lowercase letters
	var sep byte
	for {
		sep = byte(33 + rng.Intn(94))
		if sep < 'a' || sep > 'z' {
			break
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, name := range names {
		sb.WriteString(name)
		sb.WriteByte('\n')
	}
	sb.WriteByte(sep)
	sb.WriteByte('\n')
	return sb.String()
}

func randName(rng *rand.Rand, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		expected, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected:\n%s\ngot:\n%s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
