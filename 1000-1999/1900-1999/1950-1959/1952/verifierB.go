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

func run(bin, input string) (string, error) {
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

func generateCase(rng *rand.Rand) (string, []string) {
	t := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(10) + 1
		b := make([]byte, n)
		for j := 0; j < n; j++ {
			b[j] = byte('a' + rng.Intn(26))
		}
		// insert "it" with 50% probability
		if rng.Intn(2) == 0 && n >= 2 {
			pos := rng.Intn(n - 1)
			b[pos] = 'i'
			b[pos+1] = 't'
		}
		s := string(b)
		sb.WriteString(fmt.Sprintf("%s\n", s))
		if strings.Contains(s, "it") {
			expected[i] = "YES"
		} else {
			expected[i] = "NO"
		}
	}
	return sb.String(), expected
}

func check(out string, expect []string) error {
	fields := strings.Fields(out)
	if len(fields) != len(expect) {
		return fmt.Errorf("expected %d lines got %d", len(expect), len(fields))
	}
	for i, e := range expect {
		if fields[i] != e {
			return fmt.Errorf("line %d: expected %s got %s", i+1, e, fields[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if err := check(out, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
