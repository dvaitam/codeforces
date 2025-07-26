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

func solve(n, r int, pr []project) string {
	pos := make([]project, 0)
	neg := make([]project, 0)
	for _, p := range pr {
		if p.b >= 0 {
			pos = append(pos, p)
		} else {
			neg = append(neg, p)
		}
	}
	sort.Slice(pos, func(i, j int) bool {
		if pos[i].a == pos[j].a {
			return pos[i].b > pos[j].b
		}
		return pos[i].a < pos[j].a
	})
	count := 0
	curr := r
	for _, p := range pos {
		if curr >= p.a {
			curr += p.b
			count++
		}
	}
	sort.Slice(neg, func(i, j int) bool {
		left := neg[i].a + neg[i].b
		right := neg[j].a + neg[j].b
		if left == right {
			return neg[i].a > neg[j].a
		}
		return left > right
	})
	maxR := curr
	const negInf = -1000000000
	dp := make([]int, maxR+1)
	for i := range dp {
		dp[i] = negInf
	}
	dp[curr] = 0
	for _, p := range neg {
		for rating := maxR; rating >= p.a; rating-- {
			if dp[rating] == negInf {
				continue
			}
			newR := rating + p.b
			if newR < 0 {
				continue
			}
			if dp[rating]+1 > dp[newR] {
				dp[newR] = dp[rating] + 1
			}
		}
	}
	best := 0
	for _, v := range dp {
		if v > best {
			best = v
		}
	}
	return fmt.Sprintf("%d", best+count)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierF2.go /path/to/binary")
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
