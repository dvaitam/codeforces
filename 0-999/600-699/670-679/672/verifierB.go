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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expectedChanges(s string) int {
	n := len(s)
	if n > 26 {
		return -1
	}
	seen := make(map[rune]struct{})
	for _, ch := range s {
		seen[ch] = struct{}{}
	}
	return n - len(seen)
}

func randString(rng *rand.Rand, n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	sb := make([]rune, n)
	for i := 0; i < n; i++ {
		sb[i] = letters[rng.Intn(len(letters))]
	}
	return string(sb)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(100) + 1
		s := randString(rng, n)
		input := fmt.Sprintf("%d\n%s\n", n, s)
		expected := fmt.Sprintf("%d", expectedChanges(s))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
