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

type TestCase struct {
	s string
	a string
}

func encode(a string, rng *rand.Rand) string {
	var sb strings.Builder
	for i := 0; i < len(a); i++ {
		c := a[i]
		r := rng.Intn(3) // 0..2 extra letters
		for j := 0; j < r; j++ {
			ch := byte('a' + rng.Intn(26))
			for ch == c {
				ch = byte('a' + rng.Intn(26))
			}
			sb.WriteByte(ch)
		}
		sb.WriteByte(c)
	}
	return sb.String()
}

func solve(n int, s string) string {
	res := make([]byte, 0, n)
	for i := 0; i < n; {
		c := s[i]
		res = append(res, c)
		i++
		for i < n && s[i] != c {
			i++
		}
		if i < n {
			i++
		}
	}
	return string(res)
}

func generateCase(rng *rand.Rand) (string, string) {
	l := rng.Intn(20) + 1
	var ab strings.Builder
	for i := 0; i < l; i++ {
		ab.WriteByte(byte('a' + rng.Intn(26)))
	}
	a := ab.String()
	s := encode(a, rng)
	input := fmt.Sprintf("1\n%d\n%s\n", len(s), s)
	return input, a
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errb.String())
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
