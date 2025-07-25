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

func expectedAnswer(s string) string {
	b := []byte(s)
	n := len(b)
	i := 0
	for i < n && b[i] == 'a' {
		i++
	}
	if i == n {
		b[n-1] = 'z'
	} else {
		for i < n && b[i] != 'a' {
			b[i]--
			i++
		}
	}
	return string(b)
}

func runCase(bin, s string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(s + "\n")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := expectedAnswer(strings.TrimSpace(s))
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func randomString(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		s := randomString(rng)
		if err := runCase(bin, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", t+1, err, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
