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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func expected(password string) string {
	if len(password) < 5 {
		return "Too weak"
	}
	hasUpper, hasLower, hasDigit := false, false, false
	for i := 0; i < len(password); i++ {
		c := password[i]
		if c >= 'A' && c <= 'Z' {
			hasUpper = true
		} else if c >= 'a' && c <= 'z' {
			hasLower = true
		} else if c >= '0' && c <= '9' {
			hasDigit = true
		}
	}
	if hasUpper && hasLower && hasDigit {
		return "Correct"
	}
	return "Too weak"
}

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!?.,_"

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(100) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(alphabet[rng.Intn(len(alphabet))])
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin, input string) error {
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	ans := strings.TrimSpace(out)
	exp := expected(strings.TrimSpace(input))
	if ans != exp {
		return fmt.Errorf("expected %q got %q", exp, ans)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	tests := []string{
		"aaaaa\n",
		"Aaaaa\n",
		"A1a11\n",
		"AA11a\n",
		"1234\n",
	}
	for i := 0; i < 100; i++ {
		tests = append(tests, generateCase(rng))
	}

	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
