package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T, n int
	fmt.Fscan(in, &T)
	for T > 0 {
		T--
		fmt.Fscan(in, &n)
		A := make([]int, n+1)
		B := make([]int, n+1)
		parent := make([]int, n+1)
		type Ques struct{ x, id int }
		S := make([]Ques, n)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &A[i])
			B[i] = A[i]
			parent[i] = i
			S[i-1] = Ques{A[i], i}
		}
		var find func(int) int
		find = func(x int) int {
			if parent[x] != x {
				parent[x] = find(parent[x])
			}
			return parent[x]
		}
		union := func(a, b int) {
			ra := find(a)
			rb := find(b)
			if ra != rb {
				parent[ra] = rb
			}
		}
		// initial unions
		for i := 1; i <= n; i++ {
			union(i, A[i])
		}
		// sort by value
		sort.Slice(S, func(i, j int) bool {
			return S[i].x < S[j].x
		})
		// swap B for adjacent differing components
		for i := 1; i < n; i++ {
			u := S[i-1].id
			v := S[i].id
			if find(u) != find(v) {
				B[u], B[v] = B[v], B[u]
				union(u, v)
			}
		}
		// build cycle path starting from 1
		C := make([]int, 0, n)
		C = append(C, 1)
		x := B[1]
		for x != 1 {
			C = append(C, x)
			x = B[x]
		}
		// output in reverse
		for i := len(C) - 1; i >= 0; i-- {
			fmt.Fprint(out, C[i])
			if i > 0 {
				fmt.Fprint(out, " ")
			}
		}
		fmt.Fprintln(out)
	}
}
