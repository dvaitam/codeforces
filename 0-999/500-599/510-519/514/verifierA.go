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

func expected(s string) string {
	b := []byte(s)
	for i := 0; i < len(b); i++ {
		inv := '9' - b[i] + '0'
		if inv < b[i] {
			if i == 0 && inv == '0' {
				continue
			}
			b[i] = inv
		}
	}
	return string(b)
}

func genCase(rng *rand.Rand) string {
	length := rng.Intn(18) + 1
	var sb strings.Builder
	sb.WriteByte(byte(rng.Intn(9)+1) + '0')
	for i := 1; i < length; i++ {
		sb.WriteByte(byte(rng.Intn(10)) + '0')
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Deterministic cases
	cases := []string{
		"9\n", "5\n", "4\n", "4545\n", "909\n", "123456789\n", "9000\n",
		"8\n", "1111\n", "9876543210\n",
	}
	for idx, in := range cases {
		exp := expected(strings.TrimSpace(in))
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "edge case %d failed: %v\ninput:\n%s", idx+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "edge case %d failed: expected %s got %s\ninput:\n%s", idx+1, exp, out, in)
			os.Exit(1)
		}
	}

	// Random cases
	for i := 0; i < 100; i++ {
		in := genCase(rng)
		exp := expected(strings.TrimSpace(in))
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
