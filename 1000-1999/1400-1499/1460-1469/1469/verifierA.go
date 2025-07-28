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

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(99) + 2 // length between 2 and 100
	s := make([]byte, n)
	for i := range s {
		s[i] = '?'
	}
	i := rng.Intn(n)
	j := rng.Intn(n - 1)
	if j >= i {
		j++
	}
	s[i] = '('
	s[j] = ')'
	return string(s)
}

func expected(s string) string {
	n := len(s)
	if n%2 == 1 || s[0] == ')' || s[n-1] == '(' {
		return "NO"
	}
	return "YES"
}

func runCase(bin string, s string) error {
	input := fmt.Sprintf("1\n%s\n", s)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	got := strings.ToUpper(strings.TrimSpace(out))
	exp := expected(s)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s := generateCase(rng)
		if err := runCase(bin, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", i+1, err, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
