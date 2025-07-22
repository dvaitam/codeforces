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

func expected(s string) int64 {
	b := []byte(s)
	var heavyCount int64
	var result int64
	n := len(b)
	for i := 0; i+4 < n; i++ {
		if b[i] == 'h' && b[i+1] == 'e' && b[i+2] == 'a' && b[i+3] == 'v' && b[i+4] == 'y' {
			heavyCount++
		} else if b[i] == 'm' && b[i+1] == 'e' && b[i+2] == 't' && b[i+3] == 'a' && b[i+4] == 'l' {
			result += heavyCount
		}
	}
	return result
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

type test struct {
	s string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []test
	// fixed edge cases
	fixed := []string{
		"heavymetal",
		"heavyheavymetal",
		"heavymetalisheavymetal",
		"heavymetalisheavymetalisheavy",
		"metalheavy",
		"heavy",
		"metal",
		"h",
		"",
		"heavymetalm",
	}
	for _, f := range fixed {
		tests = append(tests, test{f})
	}
	for len(tests) < 100 {
		n := rng.Intn(50) + 1 // up to 50 chars for simplicity
		var b strings.Builder
		letters := []byte("heavymlt") // letters from heavy+metal
		for i := 0; i < n; i++ {
			b.WriteByte(letters[rng.Intn(len(letters))])
		}
		tests = append(tests, test{b.String()})
	}
	for i, t := range tests {
		input := t.s + "\n"
		expectedOut := fmt.Sprint(expected(t.s))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expectedOut {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expectedOut, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
