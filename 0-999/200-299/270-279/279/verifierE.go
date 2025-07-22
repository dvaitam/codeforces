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

func solveE(s string) string {
	N := len(s)
	b := make([]byte, N+2)
	for i := 0; i < N; i++ {
		b[N-1-i] = s[i] - '0'
	}
	var carry int
	var ans int64
	for i := 0; i < N+1; i++ {
		bit := int(b[i]) + carry
		if bit&1 == 0 {
			if bit == 2 {
				carry = 1
			} else {
				carry = 0
			}
		} else {
			if b[i+1] == 1 {
				carry = 1
			} else {
				carry = 0
			}
			ans++
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateCaseE(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	b := make([]byte, n)
	for i := range b {
		if rng.Intn(2) == 1 {
			b[i] = '1'
		} else {
			b[i] = '0'
		}
	}
	s := string(b)
	input := s + "\n"
	expected := solveE(s)
	return input, expected
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseE(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
