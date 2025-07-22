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

func expectedB(n int, piles []int, queries []int) []int {
	pref := make([]int, n)
	sum := 0
	for i, v := range piles {
		sum += v
		pref[i] = sum
	}
	ans := make([]int, len(queries))
	for i, q := range queries {
		for j, p := range pref {
			if q <= p {
				ans[i] = j + 1
				break
			}
		}
	}
	return ans
}

func genCase(rng *rand.Rand) (int, []int, int, []int) {
	n := rng.Intn(10) + 1
	piles := make([]int, n)
	for i := range piles {
		piles[i] = rng.Intn(20) + 1
	}
	m := rng.Intn(10) + 1
	sum := 0
	for _, v := range piles {
		sum += v
	}
	queries := make([]int, m)
	for i := range queries {
		queries[i] = rng.Intn(sum) + 1
	}
	return n, piles, m, queries
}

func runCase(bin string, n int, piles []int, m int, queries []int, exp []int) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range piles {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", m)
	for i, v := range queries {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Fields(strings.TrimSpace(out.String()))
	if len(lines) != len(exp) {
		return fmt.Errorf("expected %d lines got %d", len(exp), len(lines))
	}
	for i, v := range lines {
		var val int
		fmt.Sscan(v, &val)
		if val != exp[i] {
			return fmt.Errorf("query %d expected %d got %d", i+1, exp[i], val)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, piles, m, queries := genCase(rng)
		exp := expectedB(n, piles, queries)
		if err := runCase(bin, n, piles, m, queries, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			fmt.Fprintf(os.Stderr, "n=%d piles=%v m=%d queries=%v\n", n, piles, m, queries)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
