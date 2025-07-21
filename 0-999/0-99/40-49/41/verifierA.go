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

func reverse(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	letters := make([]byte, n)
	for i := 0; i < n; i++ {
		letters[i] = byte('a' + rng.Intn(26))
	}
	s := string(letters)
	if rng.Intn(2) == 0 {
		t := reverse(s)
		return fmt.Sprintf("%s\n%s\n", s, t), "YES\n"
	}
	// make t different
	t := reverse(s)
	idx := rng.Intn(len(t))
	tBytes := []byte(t)
	orig := tBytes[idx]
	for tBytes[idx] == orig {
		tBytes[idx] = byte('a' + rng.Intn(26))
	}
	t = string(tBytes)
	return fmt.Sprintf("%s\n%s\n", s, t), "NO\n"
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
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
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
