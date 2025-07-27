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

func solveCase(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteByte('1')
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(100) + 1
	var in strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", t))
	expLines := make([]string, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(100) + 1
		in.WriteString(fmt.Sprintf("%d\n", n))
		expLines[i] = solveCase(n)
	}
	return in.String(), strings.Join(expLines, "\n")
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	outLines := strings.FieldsFunc(strings.TrimSpace(out.String()), func(r rune) bool { return r == '\n' || r == '\r' })
	expLines := strings.FieldsFunc(strings.TrimSpace(expected), func(r rune) bool { return r == '\n' || r == '\r' })
	if len(outLines) != len(expLines) {
		return fmt.Errorf("expected %d lines got %d", len(expLines), len(outLines))
	}
	for i := range expLines {
		if strings.TrimSpace(outLines[i]) != expLines[i] {
			return fmt.Errorf("line %d expected %q got %q", i+1, expLines[i], strings.TrimSpace(outLines[i]))
		}
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
