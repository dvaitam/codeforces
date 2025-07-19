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

func expectedAnswerC(n int, s string) int {
	x := n / 2
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '(' {
			x += 2
		}
	}
	return x
}

func generateCaseC(rng *rand.Rand) (int, string) {
	n := rng.Intn(10) + 1
	b := make([]byte, n)
	for i := range b {
		if rng.Intn(2) == 0 {
			b[i] = '('
		} else {
			b[i] = ')'
		}
	}
	return n, string(b)
}

func runCaseC(bin string, n int, s string) error {
	input := fmt.Sprintf("1\n%d %s\n", n, s)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := fmt.Sprint(expectedAnswerC(n, s))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, s := generateCaseC(rng)
		if err := runCaseC(bin, n, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d %s\n", i+1, err, n, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
