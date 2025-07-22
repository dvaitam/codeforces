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

func isCapsLockError(s string) bool {
	if len(s) == 0 {
		return false
	}
	allUpper := true
	for _, c := range s {
		if c < 'A' || c > 'Z' {
			allUpper = false
			break
		}
	}
	if allUpper {
		return true
	}
	for i, c := range s {
		if i == 0 {
			continue
		}
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	return true
}

func solveCase(s string) string {
	if isCapsLockError(s) {
		r := []rune(s)
		for i, c := range r {
			if c >= 'a' && c <= 'z' {
				r[i] = c - ('a' - 'A')
			} else if c >= 'A' && c <= 'Z' {
				r[i] = c + ('a' - 'A')
			}
		}
		s = string(r)
	}
	return s
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	b := make([]rune, n)
	if rng.Float64() < 0.5 {
		// generate caps lock error
		if rng.Float64() < 0.5 {
			for i := range b {
				b[i] = rune('A' + rng.Intn(26))
			}
		} else {
			b[0] = rune('a' + rng.Intn(26))
			for i := 1; i < n; i++ {
				b[i] = rune('A' + rng.Intn(26))
			}
		}
	} else {
		for i := range b {
			if rng.Float64() < 0.5 {
				b[i] = rune('a' + rng.Intn(26))
			} else {
				b[i] = rune('A' + rng.Intn(26))
			}
		}
		// ensure not caps lock error
		if isCapsLockError(string(b)) {
			i := rng.Intn(n)
			if b[i] >= 'a' && b[i] <= 'z' {
				b[i] = b[i] - ('a' - 'A')
			} else {
				b[i] = b[i] + ('a' - 'A')
			}
		}
	}
	s := string(b)
	input := s + "\n"
	expected := solveCase(s)
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
	if outStr != expected {
		return fmt.Errorf("expected %s got %s", expected, outStr)
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
