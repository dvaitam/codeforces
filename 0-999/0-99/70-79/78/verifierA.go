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

func countVowels(s string) int {
	cnt := 0
	for _, ch := range s {
		switch ch {
		case 'a', 'e', 'i', 'o', 'u':
			cnt++
		}
	}
	return cnt
}

func generateCase(rng *rand.Rand) (string, string) {
	lines := make([]string, 3)
	counts := make([]int, 3)
	for i := 0; i < 3; i++ {
		l := rng.Intn(20) + 1
		var sb strings.Builder
		for j := 0; j < l; j++ {
			if rng.Intn(5) == 0 {
				sb.WriteByte(' ')
				continue
			}
			ch := byte('a' + rng.Intn(26))
			sb.WriteByte(ch)
			switch ch {
			case 'a', 'e', 'i', 'o', 'u':
				counts[i]++
			}
		}
		if strings.TrimSpace(sb.String()) == "" {
			sb.WriteByte('a')
			counts[i]++
		}
		lines[i] = sb.String()
	}
	input := strings.Join(lines, "\n") + "\n"
	exp := "NO"
	if counts[0] == 5 && counts[1] == 7 && counts[2] == 5 {
		exp = "YES"
	}
	return input, exp
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
		return fmt.Errorf("expected %q got %q", expected, got)
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
