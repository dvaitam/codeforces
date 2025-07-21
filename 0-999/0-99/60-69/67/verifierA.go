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

func solveA(n int, s string) string {
	candies := make([]int, n)
	for i := range candies {
		candies[i] = 1
	}
	for i := 1; i < n; i++ {
		switch s[i-1] {
		case 'R':
			candies[i] = candies[i-1] + 1
		case '=':
			candies[i] = candies[i-1]
		}
	}
	for i := n - 2; i >= 0; i-- {
		switch s[i] {
		case 'L':
			if candies[i] <= candies[i+1] {
				candies[i] = candies[i+1] + 1
			}
		case '=':
			if candies[i] != candies[i+1] {
				candies[i] = candies[i+1]
			}
		}
	}
	var sb strings.Builder
	for i, c := range candies {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", c)
	}
	return sb.String()
}

func generateCaseA(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 2
	b := make([]byte, n-1)
	letters := []byte{'L', 'R', '='}
	for i := 0; i < n-1; i++ {
		b[i] = letters[rng.Intn(len(letters))]
	}
	s := string(b)
	input := fmt.Sprintf("%d\n%s\n", n, s)
	expected := solveA(n, s)
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
	outStr := strings.TrimSpace(out.String())
	if outStr != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, outStr)
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
		in, exp := generateCaseA(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
