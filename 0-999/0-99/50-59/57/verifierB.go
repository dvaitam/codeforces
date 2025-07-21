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

func expected(n int, a, b, c []int, queries []int) int64 {
	var ans int64
	for _, q := range queries {
		for i := 0; i < len(a); i++ {
			if a[i] <= q && q <= b[i] {
				ans += int64(c[i]) + int64(q-a[i])
			}
		}
	}
	return ans
}

func genCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(50) + 1
	m := rng.Intn(20) + 1
	k := rng.Intn(n) + 1
	a := make([]int, m)
	b := make([]int, m)
	c := make([]int, m)
	for i := 0; i < m; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n) + 1
		if l > r {
			l, r = r, l
		}
		a[i] = l
		b[i] = r
		c[i] = rng.Intn(1000) + 1
	}
	queries := make([]int, k)
	used := make(map[int]bool)
	for i := 0; i < k; i++ {
		for {
			v := rng.Intn(n) + 1
			if !used[v] {
				used[v] = true
				queries[i] = v
				break
			}
		}
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for i := 0; i < m; i++ {
		fmt.Fprintf(&sb, "%d %d %d\n", a[i], b[i], c[i])
	}
	for i := 0; i < k; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", queries[i])
	}
	sb.WriteByte('\n')

	return sb.String(), expected(n, a, b, c, queries)
}

func runCase(bin, input string, exp int64) error {
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
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
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
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
