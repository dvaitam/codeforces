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
	return strings.TrimSpace(out.String()), nil
}

func encrypt(s string, k int) string {
	r := []rune(s)
	for i, ch := range r {
		if ch >= 'A' && ch <= 'Z' {
			r[i] = rune('A' + (int(ch-'A')+k)%26)
		}
	}
	return string(r)
}

func runCase(bin, s string, k int) error {
	input := fmt.Sprintf("%s %d\n", s, k)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	expected := encrypt(s, k)
	if out != expected {
		return fmt.Errorf("expected %s got %s", expected, out)
	}
	return nil
}

func randomString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('A' + rng.Intn(26))
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	edge := []struct {
		s string
		k int
	}{
		{"A", 0}, {"Z", 0}, {"ABC", 25}, {"XYZ", 1}, {"HELLOWORLD", 13},
	}
	idx := 0
	for ; idx < len(edge); idx++ {
		c := edge[idx]
		if err := runCase(bin, c.s, c.k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (s=%s k=%d)\n", idx+1, err, c.s, c.k)
			os.Exit(1)
		}
	}
	for ; idx < 100; idx++ {
		s := randomString(rng, rng.Intn(10)+1)
		k := rng.Intn(26)
		if err := runCase(bin, s, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (s=%s k=%d)\n", idx+1, err, s, k)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
