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
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case 'H', 'Q', '9':
			return "YES"
		}
	}
	return "NO"
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(100) + 1 // length 1..100
	b := make([]byte, n)
	for i := range b {
		// ASCII 33..126 inclusive
		b[i] = byte(rng.Intn(94) + 33)
	}
	return string(b)
}

func runCase(bin, input string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input + "\n")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	expected := expectedOutput(input)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
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
		input := generateCase(rng)
		if err := runCase(bin, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
