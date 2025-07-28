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

func expectedOutput(s string) string {
	b := []byte(s)
	for i := 0; i < len(b); i++ {
		if i%2 == 0 { // Alice's move (1-indexed odd)
			if b[i] == 'a' {
				b[i] = 'b'
			} else {
				b[i] = 'a'
			}
		} else { // Bob's move
			if b[i] == 'z' {
				b[i] = 'y'
			} else {
				b[i] = 'z'
			}
		}
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	s := string(b)

	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(s)
	sb.WriteByte('\n')

	expected := expectedOutput(s) + "\n"
	return sb.String(), expected
}

func runCase(exe string, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if outStr != exp {
		return fmt.Errorf("expected %q got %q", exp, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	edge := []string{
		"a",
		"z",
		"b",
		"abc",
		strings.Repeat("a", 50),
		strings.Repeat("z", 50),
		strings.Repeat("m", 50),
		"azazaz",
		"zzzza",
		"aaaaa",
	}
	for i, s := range edge {
		input := "1\n" + s + "\n"
		expected := expectedOutput(s) + "\n"
		if err := runCase(exe, input, expected); err != nil {
			fmt.Fprintf(os.Stderr, "edge case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}

	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "random case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
