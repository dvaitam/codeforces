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

func canCompose(s1, s2 string) bool {
	count := make(map[rune]int)
	for _, c := range s1 {
		if c != ' ' {
			count[c]++
		}
	}
	for _, c := range s2 {
		if c == ' ' {
			continue
		}
		if count[c] == 0 {
			return false
		}
		count[c]--
	}
	return true
}

func randString(rng *rand.Rand, n int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(5) == 0 {
			sb.WriteByte(' ')
		} else {
			sb.WriteByte(letters[rng.Intn(len(letters))])
		}
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string, string) {
	s1 := randString(rng, rng.Intn(20)+1)
	s2 := randString(rng, rng.Intn(20)+1)
	var sb strings.Builder
	sb.WriteString(s1)
	sb.WriteByte('\n')
	sb.WriteString(s2)
	sb.WriteByte('\n')
	exp := "NO"
	if canCompose(s1, s2) {
		exp = "YES"
	}
	return sb.String(), exp, fmt.Sprintf("%q\n%q\n", s1, s2)
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp, _ := generateCase(rng)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
