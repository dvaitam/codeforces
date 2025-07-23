package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int64 = -1 << 60

// Node represents a state in the Aho-Corasick automaton
type Node struct {
	next [26]int
	fail int
	val  int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var L int64
	if _, err := fmt.Fscan(in, &n, &L); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	nodes := make([]Node, 1)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		v := 0
		for j := 0; j < len(s); j++ {
			c := int(s[j] - 'a')
			if nodes[v].next[c] == 0 {
				nodes = append(nodes, Node{})
				nodes[v].next[c] = len(nodes) - 1
			}
			v = nodes[v].next[c]
		}
		nodes[v].val += a[i]
	}

	// build failure links
	queue := make([]int, 0)
	for c := 0; c < 26; c++ {
		v := nodes[0].next[c]
		if v != 0 {
			queue = append(queue, v)
		}
	}
	for qi := 0; qi < len(queue); qi++ {
		v := queue[qi]
		f := nodes[v].fail
		nodes[v].val += nodes[f].val
		for c := 0; c < 26; c++ {
			u := nodes[v].next[c]
			if u != 0 {
				nodes[u].fail = nodes[f].next[c]
				queue = append(queue, u)
			} else {
				nodes[v].next[c] = nodes[f].next[c]
			}
		}
	}

	m := len(nodes)
	// transition matrix of size m x m
	mat := make([][]int64, m)
	for i := 0; i < m; i++ {
		mat[i] = make([]int64, m)
		for j := 0; j < m; j++ {
			mat[i][j] = inf
		}
	}
	for i := 0; i < m; i++ {
		for c := 0; c < 26; c++ {
			j := nodes[i].next[c]
			val := nodes[j].val
			if val > mat[i][j] {
				mat[i][j] = val
			}
		}
	}

	// vector for current dp (only root has 0)
	dp := make([]int64, m)
	for i := 1; i < m; i++ {
		dp[i] = inf
	}

	pow := mat
	for L > 0 {
		if L&1 == 1 {
			dp = vecMul(dp, pow)
		}
		L >>= 1
		if L > 0 {
			pow = matMul(pow, pow)
		}
	}

	ans := dp[0]
	for i := 1; i < m; i++ {
		if dp[i] > ans {
			ans = dp[i]
		}
	}
	fmt.Fprintln(out, ans)
}

func vecMul(v []int64, M [][]int64) []int64 {
	n := len(v)
	res := make([]int64, n)
	for i := 0; i < n; i++ {
		res[i] = inf
	}
	for i := 0; i < n; i++ {
		if v[i] == inf {
			continue
		}
		row := M[i]
		for j := 0; j < n; j++ {
			if row[j] == inf {
				continue
			}
			val := v[i] + row[j]
			if val > res[j] {
				res[j] = val
			}
		}
	}
	return res
}

func matMul(A, B [][]int64) [][]int64 {
	n := len(A)
	C := make([][]int64, n)
	for i := 0; i < n; i++ {
		C[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			C[i][j] = inf
		}
	}
	for i := 0; i < n; i++ {
		for k := 0; k < n; k++ {
			if A[i][k] == inf {
				continue
			}
			for j := 0; j < n; j++ {
				if B[k][j] == inf {
					continue
				}
				val := A[i][k] + B[k][j]
				if val > C[i][j] {
					C[i][j] = val
				}
			}
		}
	}
	return C
}
