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

func minInt64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func solveCase(n int, a []int64, b, c []int, edges [][]int) int64 {
	var ans int64
	var dfs func(v, p int, mn int64) (int, int)
	dfs = func(v, p int, mn int64) (int, int) {
		cnt01, cnt10 := 0, 0
		if b[v] != c[v] {
			if b[v] == 0 {
				cnt01 = 1
			} else {
				cnt10 = 1
			}
		}
		for _, to := range edges[v] {
			if to == p {
				continue
			}
			x01, x10 := dfs(to, v, minInt64(mn, a[to]))
			cnt01 += x01
			cnt10 += x10
		}
		t := int(minInt64(int64(cnt01), int64(cnt10)))
		ans += 2 * int64(t) * mn
		cnt01 -= t
		cnt10 -= t
		return cnt01, cnt10
	}
	dfs01, dfs10 := dfs(0, -1, a[0])
	if dfs01 != 0 || dfs10 != 0 {
		return -1
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 1
	a := make([]int64, n)
	b := make([]int, n)
	c := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(5) + 1)
		b[i] = rng.Intn(2)
		c[i] = rng.Intn(2)
	}
	edges := make([][]int, n)
	pairs := make([][2]int, n-1)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		pairs[i-1] = [2]int{p + 1, i + 1}
		edges[p] = append(edges[p], i)
		edges[i] = append(edges[i], p)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", a[i], b[i], c[i]))
	}
	for _, e := range pairs {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	expect := solveCase(n, a, b, c, edges)
	return sb.String(), fmt.Sprintf("%d", expect)
}

func runCase(bin, in, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		if err := runCase(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
