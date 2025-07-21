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

func normalize(s string) string {
	var b strings.Builder
	prevSlash := false
	for _, ch := range s {
		if ch == '/' {
			if !prevSlash {
				b.WriteByte('/')
				prevSlash = true
			}
		} else {
			b.WriteRune(ch)
			prevSlash = false
		}
	}
	res := b.String()
	if len(res) > 1 && res[len(res)-1] == '/' {
		res = res[:len(res)-1]
	}
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	l := rng.Intn(50) + 1
	var sb strings.Builder
	sb.WriteByte('/')
	for i := 1; i < l; i++ {
		if rng.Intn(3) == 0 {
			sb.WriteByte('/')
		} else {
			sb.WriteByte(byte('a' + rng.Intn(26)))
		}
	}
	for i := 0; i < rng.Intn(3); i++ {
		sb.WriteByte('/')
	}
	input := sb.String()
	return input + "\n", normalize(input)
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
		return fmt.Errorf("expected '%s' got '%s'", expected, got)
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
