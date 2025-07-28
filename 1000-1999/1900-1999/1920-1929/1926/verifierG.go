package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const inf = int(1e9)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func expected(n int, parents []int, s string) string {
	g := make([][]int, n)
	parent := make([]int, n)
	parent[0] = -1
	for i := 1; i < n; i++ {
		a := parents[i-1]
		g[a] = append(g[a], i)
		g[i] = append(g[i], a)
	}
	order := make([]int, 0, n)
	stack := []int{0}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, to := range g[v] {
			if to == parent[v] {
				continue
			}
			parent[to] = v
			stack = append(stack, to)
		}
	}

	dp0 := make([]int, n)
	dp1 := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		v := order[i]
		var val0, val1 int
		switch s[v] {
		case 'S':
			val0 = inf
			val1 = 0
		case 'P':
			val0 = 0
			val1 = inf
		default:
			val0 = 0
			val1 = 0
		}
		for _, to := range g[v] {
			if to == parent[v] {
				continue
			}
			val0 += min(dp0[to], dp1[to]+1)
			val1 += min(dp1[to], dp0[to]+1)
		}
		dp0[v] = val0
		dp1[v] = val1
	}
	ans := min(dp0[0], dp1[0])
	return fmt.Sprint(ans)
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func randTree(rng *rand.Rand, n int) []int {
	par := make([]int, n-1)
	for i := 2; i <= n; i++ {
		par[i-2] = rng.Intn(i - 1)
	}
	return par
}

func randString(rng *rand.Rand, n int) string {
	letters := []byte{'P', 'S', 'C'}
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	type test struct {
		n   int
		par []int
		s   string
	}
	var tests []test
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 2
		par := randTree(rng, n)
		s := randString(rng, n)
		tests = append(tests, test{n: n, par: par, s: s})
	}
	for idx, t := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", t.n))
		for i, v := range t.par {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v + 1))
		}
		sb.WriteByte('\n')
		sb.WriteString(t.s)
		sb.WriteByte('\n')
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		exp := expected(t.n, t.par, t.s)
		if got != exp {
			fmt.Printf("test %d failed: expected=%s got=%s\n", idx+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
