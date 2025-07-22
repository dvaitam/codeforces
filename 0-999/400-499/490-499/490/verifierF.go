package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type caseF struct {
	n     int
	pop   []int
	edges [][2]int
	ans   int
}

func LIS(arr []int) int {
	if len(arr) == 0 {
		return 0
	}
	d := []int{arr[0]}
	for i := 1; i < len(arr); i++ {
		x := arr[i]
		l, r := 0, len(d)
		for l < r {
			m := (l + r) / 2
			if d[m] < x {
				l = m + 1
			} else {
				r = m
			}
		}
		if l == len(d) {
			d = append(d, x)
		} else {
			d[l] = x
		}
	}
	return len(d)
}

func path(adj [][]int, u, v int) []int {
	n := len(adj)
	prev := make([]int, n)
	for i := 0; i < n; i++ {
		prev[i] = -1
	}
	q := []int{u}
	prev[u] = u
	for len(q) > 0 {
		x := q[0]
		q = q[1:]
		if x == v {
			break
		}
		for _, nb := range adj[x] {
			if prev[nb] == -1 {
				prev[nb] = x
				q = append(q, nb)
			}
		}
	}
	if prev[v] == -1 {
		return nil
	}
	var p []int
	cur := v
	for cur != u {
		p = append(p, cur)
		cur = prev[cur]
	}
	p = append(p, u)
	for i, j := 0, len(p)-1; i < j; i, j = i+1, j-1 {
		p[i], p[j] = p[j], p[i]
	}
	return p
}

func solveCase(n int, pop []int, edges [][2]int) int {
	adj := make([][]int, n)
	for _, e := range edges {
		a, b := e[0]-1, e[1]-1
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	ans := 1
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			p := path(adj, i, j)
			if p == nil {
				continue
			}
			arr := make([]int, len(p))
			for k, idx := range p {
				arr[k] = pop[idx]
			}
			l := LIS(arr)
			if l > ans {
				ans = l
			}
		}
	}
	return ans
}

func genCase(r *rand.Rand, n int) caseF {
	pop := make([]int, n)
	for i := 0; i < n; i++ {
		pop[i] = r.Intn(100) + 1
	}
	edges := make([][2]int, n-1)
	for i := 1; i < n; i++ {
		p := r.Intn(i)
		edges[i-1] = [2]int{i + 1, p + 1}
	}
	ans := solveCase(n, pop, edges)
	return caseF{n, pop, edges, ans}
}

func generateTests() []caseF {
	r := rand.New(rand.NewSource(47))
	var tests []caseF
	tests = append(tests, genCase(r, 2))
	for len(tests) < 120 {
		n := r.Intn(8) + 2
		tests = append(tests, genCase(r, n))
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func verify(tc caseF, out string) error {
	v, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return fmt.Errorf("invalid integer")
	}
	if v != tc.ans {
		return fmt.Errorf("expected %d got %d", tc.ans, v)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for j, v := range tc.pop {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verify(tc, out); err != nil {
			fmt.Printf("wrong answer on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
