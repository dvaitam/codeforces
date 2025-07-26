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

func expected(names []string) string {
	counts := make([]int, 26)
	for _, name := range names {
		if len(name) > 0 {
			c := name[0] - 'a'
			if c < 26 {
				counts[c]++
			}
		}
	}
	result := 0
	for _, c := range counts {
		a := c / 2
		b := c - a
		result += a*(a-1)/2 + b*(b-1)/2
	}
	return fmt.Sprintf("%d", result)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	names := make([]string, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		l := rng.Intn(10) + 1
		b := make([]byte, l)
		for j := 0; j < l; j++ {
			b[j] = byte('a' + rng.Intn(26))
		}
		names[i] = string(b)
		sb.WriteString(names[i])
		sb.WriteByte('\n')
	}
	exp := expected(names)
	return sb.String(), exp
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
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
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
