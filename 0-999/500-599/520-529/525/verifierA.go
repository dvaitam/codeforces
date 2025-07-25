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

func runCandidate(bin, input string) (string, error) {
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

func solveCase(n int, s string) string {
	keys := make([]int, 26)
	bought := 0
	for i := 0; i < n-1; i++ {
		k := s[2*i] - 'a'
		d := s[2*i+1] - 'A'
		keys[k]++
		if keys[d] > 0 {
			keys[d]--
		} else {
			bought++
		}
	}
	return fmt.Sprintf("%d", bought)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(30) + 2
	var b strings.Builder
	for i := 0; i < (n-1)*2; i += 2 {
		key := byte('a' + rng.Intn(26))
		door := byte('A' + rng.Intn(26))
		b.WriteByte(key)
		b.WriteByte(door)
	}
	input := fmt.Sprintf("%d\n%s\n", n, b.String())
	expect := solveCase(n, b.String())
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
