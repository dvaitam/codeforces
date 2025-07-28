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

func isPalindrome(s string) bool {
	n := len(s)
	for i := 0; i < n/2; i++ {
		if s[i] != s[n-1-i] {
			return false
		}
	}
	return true
}

func generateCase(rng *rand.Rand) (string, []string) {
	t := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	exp := make([]string, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(10) + 1
		b := make([]byte, n)
		for j := 0; j < n; j++ {
			b[j] = byte('a' + rng.Intn(26))
		}
		if rng.Intn(2) == 0 {
			// make palindrome
			for j := 0; j < n/2; j++ {
				b[n-1-j] = b[j]
			}
		}
		s := string(b)
		sb.WriteString(fmt.Sprintf("%s\n", s))
		if isPalindrome(s) {
			exp[i] = "YES"
		} else {
			exp[i] = "NO"
		}
	}
	return sb.String(), exp
}

func check(out string, exp []string) error {
	fields := strings.Fields(out)
	if len(fields) != len(exp) {
		return fmt.Errorf("expected %d lines got %d", len(exp), len(fields))
	}
	for i, e := range exp {
		if fields[i] != e {
			return fmt.Errorf("line %d: expected %s got %s", i+1, e, fields[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
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
