package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type frac struct{ num, den int64 }

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
func reduce(a, b int64) (int64, int64) {
	if a == 0 && b == 0 {
		return 0, 1
	}
	if b == 0 {
		return 1, 0
	}
	if a == 0 {
		return 0, 1
	}
	g := gcd(abs(a), abs(b))
	a /= g
	b /= g
	if b < 0 {
		a = -a
		b = -b
	}
	return a, b
}
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func solveE(N int, A, B, C, D []int64) int {
	f1 := make([]frac, N)
	f2 := make([]frac, N)
	list := make([]frac, 0, 2*N)
	for i := 0; i < N; i++ {
		a1 := (A[i] + B[i]) * D[i]
		b1 := B[i] * C[i]
		a1, b1 = reduce(a1, b1)
		f1[i] = frac{a1, b1}
		list = append(list, f1[i])
		a2 := A[i] * D[i]
		b2 := B[i] * (C[i] + D[i])
		a2, b2 = reduce(a2, b2)
		f2[i] = frac{a2, b2}
		list = append(list, f2[i])
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].num != list[j].num {
			return list[i].num < list[j].num
		}
		return list[i].den < list[j].den
	})
	uniq := make([]frac, 0, len(list))
	for i, v := range list {
		if i == 0 || v != list[i-1] {
			uniq = append(uniq, v)
		}
	}
	m := make(map[frac]int, len(uniq))
	for i, v := range uniq {
		m[v] = i
	}
	V := len(uniq)
	type edge struct{ to, eid int }
	adj := make([][]edge, V)
	for i := 0; i < N; i++ {
		u := m[f1[i]]
		v := m[f2[i]]
		adj[u] = append(adj[u], edge{v, i})
		adj[v] = append(adj[v], edge{u, i})
	}
	vis := make([]bool, V)
	depth := make([]int, V)
	res := make([][2]int, 0, N/2)
	var dfs func(int, int, int) int
	dfs = func(u, pre, d int) int {
		vis[u] = true
		depth[u] = d
		cur := -1
		for _, e := range adj[u] {
			v := e.to
			var tmp int
			if vis[v] {
				tmp = -1
				if depth[u] < depth[v] {
					tmp = e.eid
				}
			} else {
				tmp = dfs(v, e.eid, d+1)
			}
			if tmp == -1 {
				continue
			}
			if cur == -1 {
				cur = tmp
			} else {
				res = append(res, [2]int{cur, tmp})
				cur = -1
			}
		}
		if cur >= 0 && pre >= 0 {
			res = append(res, [2]int{cur, pre})
			cur = -1
			pre = -1
		}
		return pre
	}
	for i := 0; i < V; i++ {
		if !vis[i] {
			dfs(i, -1, 0)
		}
	}
	return len(res)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	in := bufio.NewReader(file)
	var T int
	fmt.Fscan(in, &T)
	for idx := 1; idx <= T; idx++ {
		var N int
		fmt.Fscan(in, &N)
		A := make([]int64, N)
		B := make([]int64, N)
		C := make([]int64, N)
		D := make([]int64, N)
		for i := 0; i < N; i++ {
			fmt.Fscan(in, &A[i], &B[i], &C[i], &D[i])
		}
		expect := solveE(N, A, B, C, D)
		var input strings.Builder
		input.WriteString("1\n")
		fmt.Fprintf(&input, "%d\n", N)
		for i := 0; i < N; i++ {
			fmt.Fprintf(&input, "%d %d %d %d\n", A[i], B[i], C[i], D[i])
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", idx, err, errBuf.String())
			os.Exit(1)
		}
		fields := strings.Fields(strings.TrimSpace(out.String()))
		if len(fields) == 0 {
			fmt.Printf("case %d empty output\n", idx)
			os.Exit(1)
		}
		got, err2 := strconv.Atoi(fields[0])
		if err2 != nil || got != expect {
			fmt.Printf("case %d failed: expected %d got %s\n", idx, expect, fields[0])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", T)
}
