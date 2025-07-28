package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func rebuild(p []int, jump, cnt, block []int, b, B, n int) {
	l := b*B + 1
	r := (b + 1) * B
	if r > n {
		r = n
	}
	for i := r; i >= l; i-- {
		nxt := p[i]
		if nxt >= l && nxt <= r {
			jump[i] = jump[nxt]
			cnt[i] = cnt[nxt] + 1
		} else {
			jump[i] = nxt
			cnt[i] = 1
		}
	}
}

func solveCase(n, q int, p []int, queries [][3]int) string {
	B := 1
	for B*B < n {
		B++
	}
	block := make([]int, n+1)
	for i := 1; i <= n; i++ {
		block[i] = (i - 1) / B
	}
	jump := make([]int, n+1)
	cnt := make([]int, n+1)
	blocks := (n + B - 1) / B
	for b := 0; b < blocks; b++ {
		rebuild(p, jump, cnt, block, b, B, n)
	}
	var sb strings.Builder
	for _, qu := range queries {
		t, x, y := qu[0], qu[1], qu[2]
		if t == 1 {
			p[x], p[y] = p[y], p[x]
			rebuild(p, jump, cnt, block, block[x], B, n)
			if block[x] != block[y] {
				rebuild(p, jump, cnt, block, block[y], B, n)
			}
		} else {
			i := x
			k := y
			for k > 0 {
				if cnt[i] <= k {
					k -= cnt[i]
					i = jump[i]
				} else {
					i = p[i]
					k--
				}
			}
			fmt.Fprintf(&sb, "%d\n", i)
		}
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	q := rng.Intn(20) + 1
	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
	}
	queries := make([][3]int, q)
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			x := rng.Intn(n) + 1
			y := rng.Intn(n) + 1
			queries[i] = [3]int{1, x, y}
		} else {
			x := rng.Intn(n) + 1
			k := rng.Intn(10) + 1
			queries[i] = [3]int{2, x, k}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, q)
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&sb, "%d ", p[i])
	}
	sb.WriteByte('\n')
	for _, qu := range queries {
		fmt.Fprintf(&sb, "%d %d %d\n", qu[0], qu[1], qu[2])
	}
	return sb.String(), solveCase(n, q, p, queries)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
