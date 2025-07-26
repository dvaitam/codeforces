package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type project struct {
	a int
	b int
}

func solve(n, r int, projects []project) string {
	pos := make([]project, 0)
	neg := make([]project, 0)
	for _, p := range projects {
		if p.b >= 0 {
			pos = append(pos, p)
		} else {
			neg = append(neg, p)
		}
	}
	sort.Slice(pos, func(i, j int) bool { return pos[i].a < pos[j].a })
	for _, p := range pos {
		if r < p.a {
			return "NO"
		}
		r += p.b
	}
	sort.Slice(neg, func(i, j int) bool { return neg[i].a+neg[i].b > neg[j].a+neg[j].b })
	m := len(neg)
	dp := make([]int, m+1)
	const negInf = -1 << 30
	for i := range dp {
		dp[i] = negInf
	}
	dp[0] = r
	for i, p := range neg {
		for j := i + 1; j >= 1; j-- {
			if dp[j-1] >= p.a && dp[j-1]+p.b >= 0 {
				cand := dp[j-1] + p.b
				if cand > dp[j] {
					dp[j] = cand
				}
			}
		}
	}
	if dp[m] >= 0 {
		return "YES"
	}
	return "NO"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	r := rng.Intn(100) + 1
	pr := make([]project, n)
	for i := 0; i < n; i++ {
		a := rng.Intn(100) + 1
		b := rng.Intn(601) - 300
		pr[i] = project{a, b}
	}
	parts := make([]string, 0, 2*n+2)
	parts = append(parts, fmt.Sprintf("%d %d", n, r))
	for _, p := range pr {
		parts = append(parts, fmt.Sprintf("%d %d", p.a, p.b))
	}
	input := strings.Join(parts, "\n") + "\n"
	expect := solve(n, r, pr)
	return input, expect
}

func runCase(bin, input, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(strings.Split(out.String(), "\n")[0])
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		if err := runCase(bin, input, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
