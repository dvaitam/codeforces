package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// allSplits returns all valid (a,b,c) triples for s satisfying the lexicographic condition
func allSplits(s string) [][3]string {
	n := len(s)
	res := make([][3]string, 0)
	for i := 1; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			a := s[:i]
			b := s[i:j]
			c := s[j:]
			if (a <= b && c <= b) || (b <= a && b <= c) {
				res = append(res, [3]string{a, b, c})
			}
		}
	}
	return res
}

func verifyCase(bin string, s string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader("1\n" + s + "\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("execution error: %v", err)
	}
	line := strings.TrimSpace(string(out))
	validSplits := allSplits(s)
	if line == ":(" {
		if len(validSplits) == 0 {
			return nil
		}
		return fmt.Errorf("expected split but got :( for %s", s)
	}
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return fmt.Errorf("invalid output format for %s: %q", s, line)
	}
	a, b, c := parts[0], parts[1], parts[2]
	if a+b+c != s {
		return fmt.Errorf("concatenation mismatch for %s: got %s%s%s", s, a, b, c)
	}
	if !((a <= b && c <= b) || (b <= a && b <= c)) {
		return fmt.Errorf("lexicographic condition failed for %s", s)
	}
	found := false
	for _, v := range validSplits {
		if v[0] == a && v[1] == b && v[2] == c {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("output triple not one of valid solutions for %s", s)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := ioutil.ReadFile("problemA.txt")
	if err != nil {
		fmt.Println("failed to read problemA.txt:", err)
		os.Exit(1)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) == 0 {
		fmt.Println("problemA.txt is empty")
		os.Exit(1)
	}
	t := lines[0]
	count := 0
	fmt.Sscan(t, &count)
	if count != len(lines)-1 {
		fmt.Printf("warning: declared %d cases but got %d lines\n", count, len(lines)-1)
		if count > len(lines)-1 {
			os.Exit(1)
		}
		count = len(lines) - 1
	}
	for i := 0; i < count; i++ {
		s := strings.TrimSpace(lines[i+1])
		if err := verifyCase(bin, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
