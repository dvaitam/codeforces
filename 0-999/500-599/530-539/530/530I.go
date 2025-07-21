package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var N, K int
	fmt.Fscan(reader, &N, &K)
	adj := make([][]bool, N)
	for i := range adj {
		adj[i] = make([]bool, N)
	}
	for c := 0; c < K; c++ {
		var M int
		fmt.Fscan(reader, &M)
		idx := make([]int, M)
		for i := 0; i < M; i++ {
			fmt.Fscan(reader, &idx[i])
			idx[i]--
		}
		for i := 0; i < M; i++ {
			for j := i + 1; j < M; j++ {
				u := idx[i]
				v := idx[j]
				adj[u][v] = true
				adj[v][u] = true
			}
		}
	}
	xs := make([]int, N)
	var dfs func(int, int) bool
	dfs = func(i, C int) bool {
		if i == N {
			return true
		}
		for v := 1; v <= C; v++ {
			ok := true
			for j := 0; j < i; j++ {
				if xs[j] == v && adj[i][j] {
					ok = false
					break
				}
			}
			if !ok {
				continue
			}
			xs[i] = v
			if dfs(i+1, C) {
				return true
			}
			xs[i] = 0
		}
		return false
	}
	for C := 1; C <= N; C++ {
		for i := 0; i < N; i++ {
			xs[i] = 0
		}
		if dfs(0, C) {
			w := bufio.NewWriter(os.Stdout)
			defer w.Flush()
			for i := 0; i < N; i++ {
				if i > 0 {
					fmt.Fprint(w, " ")
				}
				fmt.Fprint(w, xs[i])
			}
			fmt.Fprintln(w)
			return
		}
	}
}
