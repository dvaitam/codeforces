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

func solveB(n, s, e int, x, a, bArr, c, d []int) string {
	used := make([]bool, n+1)
	perm := make([]int, n)
	best := int64(1<<62 - 1)
	perm[0] = s
	perm[n-1] = e
	used[s] = true
	used[e] = true
	var dfs func(pos int)
	dfs = func(pos int) {
		if pos == n-1 {
			cost := int64(0)
			last := s
			for i := 1; i < n; i++ {
				j := perm[i]
				if x[j] < x[last] {
					cost += int64(x[last]-x[j]) + int64(c[last]) + int64(bArr[j])
				} else {
					cost += int64(x[j]-x[last]) + int64(d[last]) + int64(a[j])
				}
				last = j
			}
			if cost < best {
				best = cost
			}
			return
		}
		for i := 1; i <= n; i++ {
			if used[i] || i == e {
				continue
			}
			used[i] = true
			perm[pos] = i
			dfs(pos + 1)
			used[i] = false
		}
	}
	if n == 2 {
		best = 0
		if x[e] < x[s] {
			best = int64(x[s]-x[e]) + int64(c[s]) + int64(bArr[e])
		} else {
			best = int64(x[e]-x[s]) + int64(d[s]) + int64(a[e])
		}
	} else {
		dfs(1)
	}
	return fmt.Sprintf("%d", best)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 2
	s := rng.Intn(n) + 1
	e := rng.Intn(n) + 1
	for e == s {
		e = rng.Intn(n) + 1
	}
	x := make([]int, n+1)
	cur := 0
	for i := 1; i <= n; i++ {
		cur += rng.Intn(3) + 1
		x[i] = cur
	}
	a := make([]int, n+1)
	bArr := make([]int, n+1)
	c := make([]int, n+1)
	d := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = rng.Intn(5) + 1
		bArr[i] = rng.Intn(5) + 1
		c[i] = rng.Intn(5) + 1
		d[i] = rng.Intn(5) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, s, e))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", x[i]))
	}
	sb.WriteByte('\n')
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", bArr[i]))
	}
	sb.WriteByte('\n')
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", c[i]))
	}
	sb.WriteByte('\n')
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", d[i]))
	}
	sb.WriteByte('\n')
	expected := solveB(n, s, e, x, a, bArr, c, d)
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	if outStr != expected {
		return fmt.Errorf("expected %s got %s", expected, outStr)
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
