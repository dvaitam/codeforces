package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func expected(s string) int {
	n := len(s)
	totalUpper := 0
	for i := 0; i < n; i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			totalUpper++
		}
	}
	prefixLower := 0
	prefixUpper := 0
	ans := n
	for i := 0; i <= n; i++ {
		ops := prefixLower + (totalUpper - prefixUpper)
		if ops < ans {
			ans = ops
		}
		if i < n {
			if s[i] >= 'a' && s[i] <= 'z' {
				prefixLower++
			} else {
				prefixUpper++
			}
		}
	}
	return ans
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin, s string) error {
	input := s + "\n"
	exp := expected(s)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	val, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if val != exp {
		return fmt.Errorf("expected %d got %d", exp, val)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []string{"A", "aA", "ABCdef", "abcdef", "ABCDEF"}

	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		b := make([]byte, n)
		for j := range b {
			b[j] = letters[rng.Intn(len(letters))]
		}
		cases = append(cases, string(b))
	}

	for idx, s := range cases {
		if err := runCase(bin, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", idx+1, err, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
