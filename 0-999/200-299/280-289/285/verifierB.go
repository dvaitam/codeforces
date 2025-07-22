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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n, s, t int, p []int) string {
	if s == t {
		return "0"
	}
	curr := s
	steps := 0
	for {
		curr = p[curr]
		steps++
		if curr == t {
			return fmt.Sprintf("%d", steps)
		}
		if curr == s || steps > n+5 {
			return "-1"
		}
	}
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	perm := rng.Perm(n)
	p := make([]int, n+1)
	for i, v := range perm {
		p[i+1] = v + 1
	}
	s := rng.Intn(n) + 1
	t := rng.Intn(n) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, s, t))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", p[i]))
	}
	sb.WriteByte('\n')
	expect := expected(n, s, t, p)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// some deterministic edge cases
	cases := []struct {
		n    int
		p    []int
		s, t int
		exp  string
	}{
		{1, []int{0, 1}, 1, 1, "0"},
		{3, []int{0, 2, 3, 1}, 1, 1, "0"},
		{3, []int{0, 2, 3, 1}, 1, 2, "1"},
		{3, []int{0, 2, 3, 1}, 1, 3, "2"},
		{3, []int{0, 1, 2, 3}, 1, 2, "-1"},
	}
	for _, c := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", c.n, c.s, c.t))
		for i := 1; i <= c.n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", c.p[i]))
		}
		sb.WriteByte('\n')
		out, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "edge case failed: %v\ninput:\n%s", err, sb.String())
			os.Exit(1)
		}
		if out != c.exp {
			fmt.Fprintf(os.Stderr, "edge case failed: expected %s got %s\ninput:\n%s", c.exp, out, sb.String())
			os.Exit(1)
		}
	}

	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
