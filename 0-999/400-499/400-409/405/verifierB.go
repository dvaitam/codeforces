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

type force struct {
	pos int
	dir byte
}

func compute(s string) int {
	n := len(s)
	forces := make([]force, 0, n+2)
	forces = append(forces, force{0, 'L'})
	for i := 1; i <= n; i++ {
		c := s[i-1]
		if c == 'L' || c == 'R' {
			forces = append(forces, force{i, c})
		}
	}
	forces = append(forces, force{n + 1, 'R'})
	ans := 0
	for i := 0; i+1 < len(forces); i++ {
		left := forces[i]
		right := forces[i+1]
		d := right.pos - left.pos - 1
		if d <= 0 {
			continue
		}
		if left.dir == right.dir {
			continue
		}
		if left.dir == 'L' && right.dir == 'R' {
			ans += d
		} else if left.dir == 'R' && right.dir == 'L' {
			if d%2 == 1 {
				ans++
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(30) + 1
	var b strings.Builder
	for i := 0; i < n; i++ {
		r := rng.Intn(3)
		if r == 0 {
			b.WriteByte('L')
		} else if r == 1 {
			b.WriteByte('R')
		} else {
			b.WriteByte('.')
		}
	}
	s := b.String()
	var in strings.Builder
	fmt.Fprintf(&in, "%d\n%s\n", n, s)
	out := fmt.Sprintf("%d\n", compute(s))
	return in.String(), out
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	got := strings.TrimSpace(buf.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, buf.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
