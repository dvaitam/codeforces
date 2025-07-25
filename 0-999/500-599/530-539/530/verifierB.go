package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
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

func expected(s string) string {
	n := len(s)
	first := []rune(s[:n/2])
	second := []rune(s[n/2:])
	for i, j := 0, len(first)-1; i < j; i, j = i+1, j-1 {
		first[i], first[j] = first[j], first[i]
	}
	for i, j := 0, len(second)-1; i < j; i, j = i+1, j-1 {
		second[i], second[j] = second[j], second[i]
	}
	return string(first) + string(second)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < 100; i++ {
		n := (rng.Intn(10) + 1) * 2
		runes := make([]rune, n)
		for j := 0; j < n; j++ {
			runes[j] = letters[rng.Intn(len(letters))]
		}
		s := string(runes)
		input := s + "\n"
		want := expected(s)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, want, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
