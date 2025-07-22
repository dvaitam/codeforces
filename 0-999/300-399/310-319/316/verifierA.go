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
}

func expected(s string) int64 {
	seen := make([]bool, 10)
	k := 0
	for _, c := range s {
		if c >= 'A' && c <= 'J' {
			idx := c - 'A'
			if !seen[idx] {
				seen[idx] = true
				k++
			}
		}
	}
	var letterMapping int64 = 1
	first := s[0]
	if first >= 'A' && first <= 'J' {
		letterMapping = 9
		rem := k - 1
		for i := 0; i < rem; i++ {
			letterMapping *= int64(9 - i)
		}
	} else {
		for i := 0; i < k; i++ {
			letterMapping *= int64(10 - i)
		}
	}
	multiplier := int64(1)
	for i, c := range s {
		if c == '?' {
			if i == 0 {
				multiplier *= 9
			} else {
				multiplier *= 10
			}
		}
	}
	return letterMapping * multiplier
}

func genCase(rng *rand.Rand) (string, string) {
	length := rng.Intn(10) + 1
	firstChoices := []byte("123456789?ABCDEFGHIJ")
	otherChoices := []byte("0123456789?ABCDEFGHIJ")
	bs := make([]byte, length)
	bs[0] = firstChoices[rng.Intn(len(firstChoices))]
	for i := 1; i < length; i++ {
		bs[i] = otherChoices[rng.Intn(len(otherChoices))]
	}
	s := string(bs)
	input := fmt.Sprintf("%s\n", s)
	output := fmt.Sprintf("%d\n", expected(s))
	return input, output
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
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
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
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
