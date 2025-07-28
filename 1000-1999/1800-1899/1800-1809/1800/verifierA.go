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

func isMeow(s string) bool {
	s = strings.ToLower(s)
	letters := []byte{'m', 'e', 'o', 'w'}
	idx := 0
	for _, ch := range letters {
		if idx >= len(s) || s[idx] != ch {
			return false
		}
		for idx < len(s) && s[idx] == ch {
			idx++
		}
	}
	return idx == len(s)
}

func solveA(n int, s string) string {
	if isMeow(s) {
		return "YES"
	}
	return "NO"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(letters[rng.Intn(len(letters))])
	}
	s := sb.String()
	input := fmt.Sprintf("1\n%d\n%s\n", n, s)
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
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
