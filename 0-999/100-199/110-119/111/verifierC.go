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

func solveC(n, m int) int {
	if n > m {
		n, m = m, n
	}
	switch n {
	case 1, 2:
		return n*m - (m+2)/(4-n)
	case 3:
		m0 := m
		res := m0*n - (m0/4)*3
		r := m0 % 4
		res -= r
		if r == 0 {
			res--
		}
		return res
	case 4:
		res := m*n - m
		if m == 5 || m == 6 || m == 9 {
			res--
		}
		return res
	case 5:
		m0 := m
		res := m0*n - (m0/5)*6
		if m0 == 7 {
			res++
		}
		r := m0 % 5
		res -= r
		res--
		if r > 1 {
			res--
		}
		return res
	case 6:
		return m*n - 10
	default:
		return n * m
	}
}

func generateCase(rng *rand.Rand) (string, string) {
	var n, m int
	for {
		n = rng.Intn(8) + 1
		m = rng.Intn(8) + 1
		if n*m <= 40 {
			break
		}
	}
	input := fmt.Sprintf("%d %d\n", n, m)
	out := fmt.Sprintf("%d\n", solveC(n, m))
	return input, out
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expected), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
