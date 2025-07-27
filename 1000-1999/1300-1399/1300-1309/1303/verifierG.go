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

func solveCase(n int, edges [][2]int, a []int64) string {
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	par := make([]int, n)
	order := make([]int, 0, n)
	stack := []int{0}
	par[0] = -1
	for i := 0; i < len(stack); i++ {
		u := stack[i]
		order = append(order, u)
		for _, v := range adj[u] {
			if v == par[u] {
				continue
			}
			par[v] = u
			stack = append(stack, v)
		}
	}
	dp1 := make([]int64, n)
	sumA1 := make([]int64, n)
	for i := n - 1; i >= 0; i-- {
		u := order[i]
		dp1[u] = a[u]
		sumA1[u] = a[u]
		var bestDir, bestSum int64
		for _, v := range adj[u] {
			if v == par[u] {
				continue
			}
			dir := dp1[v] + sumA1[v]
			if dir > bestDir {
				bestDir = dir
				bestSum = sumA1[v]
			}
		}
		if bestDir > 0 {
			dp1[u] = a[u] + bestDir
			sumA1[u] = a[u] + bestSum
		}
	}
	dp2 := make([]int64, n)
	sumA2 := make([]int64, n)
	const inf = int64(4e18)
	for _, u := range order {
		var max1Val, max2Val int64 = -inf, -inf
		var max1Sum, max2Sum int64
		var max1Id int = -1
		if par[u] != -1 {
			d := dp2[u]
			s := sumA2[u]
			val := d + s
			if val > max1Val {
				max2Val, max2Sum = max1Val, max1Sum
				max1Val, max1Sum, max1Id = val, s, par[u]
			} else if val > max2Val {
				max2Val, max2Sum = val, s
			}
		}
		for _, v := range adj[u] {
			if v == par[u] {
				continue
			}
			s := sumA1[v]
			d := dp1[v] + s
			val := d + s
			if val > max1Val {
				max2Val, max2Sum = max1Val, max1Sum
				max1Val, max1Sum, max1Id = val, s, v
			} else if val > max2Val {
				max2Val, max2Sum = val, s
			}
		}
		for _, v := range adj[u] {
			if v == par[u] {
				continue
			}
			var bestVal, bestSum int64
			if max1Id != v {
				bestVal = max1Val
				bestSum = max1Sum
			} else {
				bestVal = max2Val
				bestSum = max2Sum
			}
			if bestVal > 0 {
				dp2[v] = 2*a[u] + bestVal
				sumA2[v] = a[u] + bestSum
			} else {
				dp2[v] = 0
				sumA2[v] = 0
			}
		}
	}
	var ans int64
	for i := 0; i < n; i++ {
		f1 := dp1[i]
		f2 := a[i] + dp2[i]
		if f1 > ans {
			ans = f1
		}
		if f2 > ans {
			ans = f2
		}
	}
	return fmt.Sprintf("%d\n", ans)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 2
	edges := make([][2]int, n-1)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		edges[i-1] = [2]int{p, i}
	}
	a := make([]int64, n)
	for i := range a {
		a[i] = int64(rng.Intn(10) + 1)
	}
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "%d\n", n)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0]+1, e[1]+1)
	}
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	out := solveCase(n, edges, a)
	return sb.String(), out
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	got := strings.TrimSpace(buf.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected \n%s\ngot \n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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
