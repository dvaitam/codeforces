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

func solveCase(n, k int, s string) string {
	best := 1
	for i := 1; i < n; i++ {
		if s[i] > s[i%best] {
			break
		}
		if s[i] < s[i%best] {
			best = i + 1
		}
	}
	prefix := s[:best]
	res := make([]byte, k)
	for i := 0; i < k; i++ {
		res[i] = prefix[i%best]
	}
	return string(res) + "\n"
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
	return out.String(), nil
}

func randomString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(3))
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	k := rng.Intn(20) + 1
	s := randomString(rng, n)
	input := fmt.Sprintf("%d %d\n%s\n", n, k, s)
	expect := solveCase(n, k, s)
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/binary")
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
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), strings.TrimSpace(out), in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
