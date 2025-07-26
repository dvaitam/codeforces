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

type pair struct {
	win int
	sum int64
}

var (
	n, m int
	diff []int64
	g    [][]int
)

func dfs(u, p int) ([]pair, int) {
	dp := make([]pair, 2)
	dp[1] = pair{0, diff[u]}
	size := 1
	for _, v := range g[u] {
		if v == p {
			continue
		}
		child, sz := dfs(v, u)
		newSize := size + sz
		if newSize > m {
			newSize = m
		}
		ndp := make([]pair, newSize+1)
		for i := 1; i <= size && i <= m; i++ {
			for j := 1; j <= sz && i+j-1 <= m; j++ {
				w1 := dp[i].win + child[j].win
				s1 := dp[i].sum + child[j].sum
				if w1 > ndp[i+j-1].win || (w1 == ndp[i+j-1].win && s1 > ndp[i+j-1].sum) {
					ndp[i+j-1] = pair{w1, s1}
				}
				w2 := dp[i].win + child[j].win
				if child[j].sum > 0 {
					w2++
				}
				s2 := dp[i].sum
				if w2 > ndp[i+j].win || (w2 == ndp[i+j].win && s2 > ndp[i+j].sum) {
					ndp[i+j] = pair{w2, s2}
				}
			}
		}
		size = newSize
		dp = ndp
	}
	return dp, size
}

func solveD(nv, mv int, bees, wasps []int64, edges [][2]int) int {
	n = nv
	m = mv
	diff = make([]int64, n+1)
	for i := 1; i <= n; i++ {
		diff[i] = wasps[i-1] - bees[i-1]
	}
	g = make([][]int, n+1)
	for _, e := range edges {
		x, y := e[0], e[1]
		g[x] = append(g[x], y)
		g[y] = append(g[y], x)
	}
	dp, _ := dfs(1, 0)
	ans := dp[m].win
	if dp[m].sum > 0 {
		ans++
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	t := 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(5) + 2
		m := rand.Intn(n-1) + 1
		bees := make([]int64, n)
		wasps := make([]int64, n)
		for j := 0; j < n; j++ {
			bees[j] = int64(rand.Intn(5))
			wasps[j] = int64(rand.Intn(5))
		}
		edges := make([][2]int, n-1)
		for j := 2; j <= n; j++ {
			p := rand.Intn(j-1) + 1
			edges[j-2] = [2]int{p, j}
		}
		fmt.Fprintln(&input, n, m)
		for _, v := range bees {
			fmt.Fprint(&input, v, " ")
		}
		fmt.Fprintln(&input)
		for _, v := range wasps {
			fmt.Fprint(&input, v, " ")
		}
		fmt.Fprintln(&input)
		for _, e := range edges {
			fmt.Fprintln(&input, e[0], e[1])
		}
		ans := solveD(n, m, bees, wasps, edges)
		expected[i] = fmt.Sprintf("%d", ans)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = &input
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("binary error:", err)
		fmt.Print(string(out))
		return
	}
	outputs := strings.Fields(string(out))
	if len(outputs) != t {
		fmt.Printf("expected %d lines, got %d\n", t, len(outputs))
		fmt.Print(string(out))
		return
	}
	for i := 0; i < t; i++ {
		if outputs[i] != expected[i] {
			fmt.Printf("Test %d failed: expected %s got %s\n", i+1, expected[i], outputs[i])
			return
		}
	}
	fmt.Println("All tests passed!")
}
