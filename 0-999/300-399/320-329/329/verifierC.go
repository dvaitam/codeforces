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

// We reuse the algorithm from 329C.go

func solveC(n, m int, edges [][2]int) string {
	b := make([][]int, n+1)
	cl := make([]int, n)
	leftCnt := make([]int, n+1)
	bns := make([][2]int, m)
	ans := make([][2]int, m)
	for _, e := range edges {
		u := e[0]
		v := e[1]
		b[u] = append(b[u], v)
		b[v] = append(b[v], u)
	}
	for i := 0; i < n; i++ {
		cl[i] = i + 1
	}
	var flagOK bool
	var f func(x, y, z int) bool
	f = func(x, y, z int) bool {
		for _, v := range b[x] {
			if v == y || v == z {
				return true
			}
		}
		for _, v := range b[z] {
			if v == y {
				return true
			}
		}
		return false
	}
	var btrk func(int, int, int, int, int)
	btrk = func(ln, lm, loc, p, q int) {
		if loc == lm {
			flagOK = true
			return
		}
		if p == ln {
			return
		}
		if p == q {
			btrk(ln, lm, loc, p+1, 0)
			return
		}
		u := cl[p]
		v := cl[q]
		if u > 0 && v > 0 && u != v && leftCnt[u] < 2 && leftCnt[v] < 2 {
			ok := true
			for _, w := range b[u] {
				if w == v {
					ok = false
					break
				}
			}
			if ok {
				bns[loc][0] = u
				bns[loc][1] = v
				leftCnt[u]++
				leftCnt[v]++
				btrk(ln, lm, loc+1, p, q+1)
				leftCnt[u]--
				leftCnt[v]--
				if flagOK {
					return
				}
			}
		}
		btrk(ln, lm, loc, p, q+1)
	}
	if n <= 7 {
		btrk(n, m, 0, 0, 0)
		if !flagOK {
			return "-1"
		}
		var sb strings.Builder
		for i := 0; i < m; i++ {
			fmt.Fprintf(&sb, "%d %d\n", bns[i][0], bns[i][1])
		}
		return strings.TrimSpace(sb.String())
	}
	var i int
	for i = 0; i < n-7; i += 3 {
		var jj, kk, ll int
		found := false
		for j := 0; j < 7 && !found; j++ {
			for k := 0; k < j && !found; k++ {
				for l := 0; l < k; l++ {
					if !f(cl[n-i-1-j], cl[n-i-1-k], cl[n-i-1-l]) {
						jj, kk, ll = j, k, l
						found = true
						break
					}
				}
			}
		}
		ans[i][0] = cl[n-i-1-jj]
		ans[i][1] = cl[n-i-1-kk]
		ans[i+1][0] = cl[n-i-1-jj]
		ans[i+1][1] = cl[n-i-1-ll]
		ans[i+2][0] = cl[n-i-1-kk]
		ans[i+2][1] = cl[n-i-1-ll]
		cl[n-i-1-jj] = 0
		cl[n-i-1-kk] = 0
		cl[n-i-1-ll] = 0
		for t := 0; t < 3; t++ {
			idx := n - i - 1 - t
			if cl[idx] != 0 {
				pos := n - i - 4
				for pos >= 0 && cl[pos] != 0 {
					pos--
				}
				if pos >= 0 {
					cl[pos] = cl[idx]
				}
			}
		}
	}
	var sb strings.Builder
	if i >= m {
		for j := 0; j < m; j++ {
			fmt.Fprintf(&sb, "%d %d\n", ans[j][0], ans[j][1])
		}
		return strings.TrimSpace(sb.String())
	}
	for j := 0; j < i; j++ {
		fmt.Fprintf(&sb, "%d %d\n", ans[j][0], ans[j][1])
	}
	btrk(n-i, m-i, 0, 0, 0)
	if !flagOK {
		return "-1"
	}
	for j := 0; j < m-i; j++ {
		fmt.Fprintf(&sb, "%d %d\n", bns[j][0], bns[j][1])
	}
	return strings.TrimSpace(sb.String())
}

func generateCaseC(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	m := rng.Intn(n) + 1
	deg := make([]int, n+1)
	edges := make([][2]int, 0, m)
	used := make(map[[2]int]bool)
	attempts := 0
	for len(edges) < m && attempts < 1000 {
		attempts++
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v || deg[u] >= 2 || deg[v] >= 2 {
			continue
		}
		if u > v {
			u, v = v, u
		}
		p := [2]int{u, v}
		if used[p] {
			continue
		}
		used[p] = true
		deg[u]++
		deg[v]++
		edges = append(edges, [2]int{u, v})
	}
	m = len(edges)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	out := solveC(n, m, edges)
	return sb.String(), out
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseC(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
