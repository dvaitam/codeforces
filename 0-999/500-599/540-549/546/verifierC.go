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

func solve(n int, p1, p2 []int) string {
	seen := make(map[string]struct{})
	rounds := 0
	for {
		state := encode(p1, p2)
		if _, ok := seen[state]; ok {
			return "-1"
		}
		seen[state] = struct{}{}
		if len(p1) == 0 {
			return fmt.Sprintf("%d 2", rounds)
		}
		if len(p2) == 0 {
			return fmt.Sprintf("%d 1", rounds)
		}
		rounds++
		c1 := p1[0]
		c2 := p2[0]
		p1 = p1[1:]
		p2 = p2[1:]
		if c1 > c2 {
			p1 = append(p1, c2, c1)
		} else {
			p2 = append(p2, c1, c2)
		}
	}
}

func encode(a, b []int) string {
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(',')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('|')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(',')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(9) + 2
	perm := rng.Perm(n)
	k1 := rng.Intn(n-1) + 1
	k2 := n - k1
	p1 := make([]int, k1)
	p2 := make([]int, k2)
	copy(p1, perm[:k1])
	copy(p2, perm[k1:])
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	sb.WriteString(fmt.Sprintf("%d ", k1))
	for i, v := range p1 {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v+1)
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d ", k2))
	for i, v := range p2 {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v+1)
	}
	sb.WriteByte('\n')
	input := sb.String()
	// convert to 1-based deck values for solver as well
	p1sol := make([]int, k1)
	p2sol := make([]int, k2)
	for i := range p1 {
		p1sol[i] = p1[i] + 1
	}
	for i := range p2 {
		p2sol[i] = p2[i] + 1
	}
	exp := solve(n, p1sol, p2sol)
	return input, exp
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
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
