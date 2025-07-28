package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
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
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

func insertOneA(s, t string) bool {
	if len(t) != len(s)+1 {
		return false
	}
	for k := 0; k < len(t); k++ {
		if t[k] == 'a' {
			if t[:k]+t[k+1:] == s {
				return true
			}
		}
	}
	return false
}

func validateOutput(s, out string) bool {
	out = strings.TrimSpace(out)
	if out == "NO" {
		for i := 0; i < len(s); i++ {
			if s[i] != 'a' {
				return false
			}
		}
		return true
	}
	parts := strings.Split(strings.ReplaceAll(out, "\r", ""), "\n")
	if len(parts) != 2 || strings.ToUpper(strings.TrimSpace(parts[0])) != "YES" {
		return false
	}
	t := strings.TrimSpace(parts[1])
	if !insertOneA(s, t) {
		return false
	}
	if isPalindrome(t) {
		return false
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := []string{"a", "aa", "b", "ab", "aba", "bbb", "caa", "abc", "zz", "aba", "aaaa", "ba", "cab"}
	rng := rand.New(rand.NewSource(42))
	letters := "abc"
	for len(tests) < 100 {
		n := rng.Intn(10) + 1
		var sb strings.Builder
		for i := 0; i < n; i++ {
			sb.WriteByte(letters[rng.Intn(len(letters))])
		}
		tests = append(tests, sb.String())
	}
	for idx, s := range tests {
		input := fmt.Sprintf("1\n%s\n", s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if !validateOutput(s, got) {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid output for input %q\nGot:\n%s\n", idx+1, s, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
