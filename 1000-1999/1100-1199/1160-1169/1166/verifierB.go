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

func expected(n int) string {
	if n < 25 {
		return "-1"
	}
	y := 0
	for i := 5; i*i <= n; i++ {
		if n%i == 0 {
			y = i
			break
		}
	}
	if y == 0 {
		return "-1"
	}
	other := n / y
	if other < 5 {
		return "-1"
	}
	if y == 5 {
		patterns := []string{"aeiou", "eioua", "iouae", "ouaei", "uaeio"}
		var b strings.Builder
		for i := 0; i < other; i++ {
			b.WriteString(patterns[i%5])
		}
		return b.String()
	}
	vowels := []byte{'a', 'e', 'i', 'o', 'u'}
	var b strings.Builder
	x := 0
	for j := 0; j < other; j++ {
		for i := 0; i < y; i++ {
			b.WriteByte(vowels[x%5])
			x++
		}
	}
	return b.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(200) + 1
	return fmt.Sprintf("%d\n", n), expected(n)
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
