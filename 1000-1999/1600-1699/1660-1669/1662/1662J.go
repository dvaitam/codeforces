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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	A := make([][]int, n)
	for i := 0; i < n; i++ {
		A[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &A[i][j])
		}
	}

	C := make([][]int, n)
	for i := 0; i < n; i++ {
		C[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &C[i][j])
		}
	}

	// positions of each age in rows and columns
	posRow := make([][]int, n+1)
	posCol := make([][]int, n+1)
	for age := 1; age <= n; age++ {
		posRow[age] = make([]int, n)
		posCol[age] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			val := A[i][j]
			posRow[val][i] = j
			posCol[val][j] = i
		}
	}

	best := 0
	// case: all kids of same age
	for age := 1; age <= n; age++ {
		cnt := 0
		for i := 0; i < n; i++ {
			cnt += C[i][posRow[age][i]]
		}
		if cnt > best {
			best = cnt
		}
	}

	// case: mix of two consecutive ages
	for age := 1; age < n; age++ {
		pk := posRow[age]
		pk1 := posRow[age+1]
		rowFromColPk1 := posCol[age+1]

		visited := make([]bool, n)
		total := 0
		for r := 0; r < n; r++ {
			if visited[r] {
				continue
			}
			cycle := []int{}
			x := r
			for !visited[x] {
				visited[x] = true
				cycle = append(cycle, x)
				x = rowFromColPk1[pk[x]]
			}
			gk := 0
			gk1 := 0
			for _, v := range cycle {
				gk += C[v][pk[v]]
				gk1 += C[v][pk1[v]]
			}
			if gk > gk1 {
				total += gk
			} else {
				total += gk1
			}
		}
		if total > best {
			best = total
		}
	}

	fmt.Fprintln(out, best)
}
