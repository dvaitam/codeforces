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

func isPal(s string) bool {
	i, j := 0, len(s)-1
	for i < j {
		if s[i] != s[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func possibleInsertions(s string) []string {
	results := []string{}
	for pos := 0; pos <= len(s); pos++ {
		for ch := 'a'; ch <= 'z'; ch++ {
			t := s[:pos] + string(ch) + s[pos:]
			if isPal(t) {
				results = append(results, t)
			}
		}
	}
	return results
}

func verifyOutput(s, out string) bool {
	out = strings.TrimSpace(out)
	if out == "NA" {
		return len(possibleInsertions(s)) == 0
	}
	if len(out) != len(s)+1 {
		return false
	}
	if !isPal(out) {
		return false
	}
	for i := 0; i <= len(s); i++ {
		if out[:i]+out[i+1:] == s {
			return true
		}
	}
	return false
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		l := rng.Intn(10) + 1
		b := make([]byte, l)
		for j := 0; j < l; j++ {
			b[j] = byte('a' + rng.Intn(26))
		}
		s := string(b)
		input := s + "\n"
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if !verifyOutput(s, out) {
			fmt.Fprintf(os.Stderr, "case %d failed: wrong output %s\ninput:%s", i+1, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
