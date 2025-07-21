package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		size := n * m
		rowA := make([]int, size+1)
		colA := make([]int, size+1)
		rowB := make([]int, size+1)
		colB := make([]int, size+1)
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				var x int
				fmt.Fscan(in, &x)
				rowA[x] = i
				colA[x] = j
			}
		}
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				var x int
				fmt.Fscan(in, &x)
				rowB[x] = i
				colB[x] = j
			}
		}
		rowMap := make([]int, n)
		colMap := make([]int, m)
		for i := range rowMap {
			rowMap[i] = -1
		}
		for i := range colMap {
			colMap[i] = -1
		}
		ok := true
		for val := 1; val <= size && ok; val++ {
			rA, cA := rowA[val], colA[val]
			rB, cB := rowB[val], colB[val]
			if rowMap[rA] == -1 {
				rowMap[rA] = rB
			} else if rowMap[rA] != rB {
				ok = false
				break
			}
			if colMap[cA] == -1 {
				colMap[cA] = cB
			} else if colMap[cA] != cB {
				ok = false
				break
			}
		}
		usedRow := make([]bool, n)
		usedCol := make([]bool, m)
		for _, v := range rowMap {
			if v == -1 || usedRow[v] {
				ok = false
				break
			}
			usedRow[v] = true
		}
		for _, v := range colMap {
			if v == -1 || usedCol[v] {
				ok = false
				break
			}
			usedCol[v] = true
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
