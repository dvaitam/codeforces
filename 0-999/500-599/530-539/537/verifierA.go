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

const targetA = "CODEFORCES"

func expectedAnswerA(s string) string {
	n := len(s)
	m := len(targetA)
	if n < m {
		return "NO"
	}
	for i := 0; i <= m; i++ {
		if s[:i] == targetA[:i] && s[n-(m-i):] == targetA[i:] {
			return "YES"
		}
	}
	return "NO"
}

func generateCaseA(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('A' + rng.Intn(26))
	}
	if string(b) == targetA {
		b[0] = 'Z'
	}
	return string(b)
}

func runCaseA(bin, s string) error {
	input := s + "\n"
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := expectedAnswerA(s)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s := generateCaseA(rng)
		if err := runCaseA(bin, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", i+1, err, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
