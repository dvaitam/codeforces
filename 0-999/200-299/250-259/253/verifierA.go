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

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func transitions(s string) int {
	t := 0
	for i := 0; i+1 < len(s); i++ {
		if s[i] != s[i+1] {
			t++
		}
	}
	return t
}

func maxTransitions(n, m int) int {
	if n == m {
		return 2*n - 1
	}
	if n < m {
		n, m = m, n
	}
	// now n > m
	return 2 * m
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 1; tc <= 100; tc++ {
		n := rng.Intn(100) + 1
		m := rng.Intn(100) + 1
		input := fmt.Sprintf("%d %d\n", n, m)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tc, err, input)
			os.Exit(1)
		}
		s := strings.TrimSpace(out)
		if len(s) != n+m {
			fmt.Fprintf(os.Stderr, "case %d failed: expected length %d got %d\ninput:\n%s", tc, n+m, len(s), input)
			os.Exit(1)
		}
		countB := strings.Count(s, "B")
		countG := strings.Count(s, "G")
		if countB != n || countG != m || countB+countG != len(s) {
			fmt.Fprintf(os.Stderr, "case %d failed: wrong counts B=%d G=%d output=%s\ninput:\n%s", tc, countB, countG, s, input)
			os.Exit(1)
		}
		texp := maxTransitions(n, m)
		if transitions(s) != texp {
			fmt.Fprintf(os.Stderr, "case %d failed: transitions %d expected %d output=%s\ninput:\n%s", tc, transitions(s), texp, s, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
