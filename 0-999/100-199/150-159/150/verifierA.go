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

// solve returns the expected output for problem A given a
func solve(a int64) string {
	if a == 1 {
		return "1\n0"
	}
	factors := make([]int64, 0, 3)
	n := a
	for n%2 == 0 {
		factors = append(factors, 2)
		if len(factors) > 2 {
			break
		}
		n /= 2
	}
	for i := int64(3); i*i <= n && len(factors) <= 2; i += 2 {
		for n%i == 0 {
			factors = append(factors, i)
			if len(factors) > 2 {
				break
			}
			n /= i
		}
	}
	if len(factors) <= 2 && n > 1 {
		factors = append(factors, n)
	}
	if len(factors) == 1 {
		return "1\n0"
	}
	if len(factors) == 2 && factors[0]*factors[1] == a {
		return "2"
	}
	return fmt.Sprintf("1\n%d", factors[0]*factors[1])
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		a := rng.Int63n(1000000) + 1
		input := fmt.Sprintf("%d\n", a)
		expected := solve(a)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(expected) != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
