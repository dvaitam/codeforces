package main

import (
	"bufio"
	"fmt"
	"os"
)

const N = 100010

func dfs(u int, G [][]int, col []int, ok *int) {
	for _, v := range G[u] {
		if col[v] != -1 {
			if col[u] == col[v] {
				*ok = 0
			}
		} else {
			col[v] = col[u] ^ 1
			dfs(v, G, col, ok)
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		G := make([][]int, n+1)
		col := make([]int, n+1)
		for i := 1; i <= n; i++ {
			col[i] = -1
		}
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			G[u] = append(G[u], v)
			G[v] = append(G[v], u)
		}
		col[1] = 0
		ok := 1
		dfs(1, G, col, &ok)
		if ok == 0 {
			fmt.Println("Alice")
			for i := 1; i <= n; i++ {
				fmt.Println(1, 2)
				var x, y int
				fmt.Fscan(in, &x, &y)
			}
		} else {
			fmt.Println("Bob")
			A := []int{}
			B := []int{}
			for i := 1; i <= n; i++ {
				if col[i] != 0 {
					A = append(A, i)
				} else {
					B = append(B, i)
				}
			}
			for i := 1; i <= n; i++ {
				var x, y int
				fmt.Fscan(in, &x, &y)
				if (x == 1 || y == 1) && len(A) > 0 {
					fmt.Println(A[len(A)-1], 1)
					A = A[:len(A)-1]
				} else if (x == 2 || y == 2) && len(B) > 0 {
					fmt.Println(B[len(B)-1], 2)
					B = B[:len(B)-1]
				} else if len(A) > 0 {
					fmt.Println(A[len(A)-1], 3)
					A = A[:len(A)-1]
				} else {
					fmt.Println(B[len(B)-1], 3)
					B = B[:len(B)-1]
				}
			}
		}
	}
}
