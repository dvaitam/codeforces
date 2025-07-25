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

func solveB(s string) string {
	n := len(s)
	last := s[n-1]
	earliest := -1
	latestGreater := -1
	for i := 0; i < n-1; i++ {
		c := s[i]
		if (c-'0')%2 == 0 {
			if c < last && earliest == -1 {
				earliest = i
			}
			if c > last {
				latestGreater = i
			}
		}
	}
	idx := -1
	if earliest != -1 {
		idx = earliest
	} else if latestGreater != -1 {
		idx = latestGreater
	} else {
		return "-1"
	}
	b := []byte(s)
	b[idx], b[n-1] = b[n-1], b[idx]
	if b[0] == '0' {
		return "-1"
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, string) {
	length := rng.Intn(18) + 2
	digits := make([]byte, length)
	digits[0] = byte(rng.Intn(9)+1) + '0'
	for i := 1; i < length-1; i++ {
		digits[i] = byte(rng.Intn(10)) + '0'
	}
	odd := []byte{'1', '3', '5', '7', '9'}
	digits[length-1] = odd[rng.Intn(len(odd))]
	s := string(digits)
	expect := solveB(s)
	return s + "\n", expect
}

func runCase(bin string, input string, expected string) error {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		if err := runCase(bin, input, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
