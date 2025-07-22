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
	best := ""
	bestCount := 0
	for i := 0; i < n; i++ {
		if s[i] != '4' && s[i] != '7' {
			continue
		}
		for j := i; j < n && (s[j] == '4' || s[j] == '7'); j++ {
			substr := s[i : j+1]
			count := 0
			for k := 0; k+len(substr) <= n; k++ {
				if s[k:k+len(substr)] == substr {
					count++
				}
			}
			if count > bestCount || (count == bestCount && (best == "" || substr < best)) {
				best = substr
				bestCount = count
			}
		}
	}
	if bestCount == 0 {
		return "-1"
	}
	return best
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	special := []string{"4", "7", "47", "444", "123", "7777777", "123447744"}
	for i, s := range special {
		exp := solveB(s)
		input := fmt.Sprintf("%s\n", s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("special case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("special case %d failed: s=%s expected %s got %s\n", i+1, s, exp, got)
			os.Exit(1)
		}
	}

	for i := 0; i < 100; i++ {
		length := rng.Intn(50) + 1
		b := make([]byte, length)
		for j := 0; j < length; j++ {
			b[j] = byte('0' + rng.Intn(10))
		}
		s := string(b)
		exp := solveB(s)
		input := fmt.Sprintf("%s\n", s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed: s=%s expected %s got %s\n", i+1, s, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
