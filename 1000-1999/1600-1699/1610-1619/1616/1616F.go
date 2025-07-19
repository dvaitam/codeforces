package main

import (
	"bufio"
	"fmt"
	"os"
)

var rdr = bufio.NewReader(os.Stdin)
var wrtr = bufio.NewWriter(os.Stdout)

func invThree(x int8) int8 {
	// In modulo 3, multiplicative inverse of 1 is 1, of 2 is 2 (since 2*2=4≡1)
	return x
}

func solveOne() {
	var n, m int
	if _, err := fmt.Fscan(rdr, &n, &m); err != nil {
		return
	}
	// adjacency matrix for edge indices
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		adj[i] = make([]int, n)
		for j := 0; j < n; j++ {
			adj[i][j] = -1
		}
	}
	edges := make([][2]int, m)
	cVal := make([]int, m)
	for i := 0; i < m; i++ {
		var a, b, c int
		fmt.Fscan(rdr, &a, &b, &c)
		a--
		b--
		edges[i] = [2]int{a, b}
		cVal[i] = c
		adj[a][b] = i
		adj[b][a] = i
	}
	// build system of equations mod 3
	eqA := make([][]int8, 0)
	values := make([]int8, 0)
	addEq := func(row []int8, rhs int8) {
		eqA = append(eqA, row)
		values = append(values, rhs)
	}
	// known edge weights
	for i := 0; i < m; i++ {
		if cVal[i] != -1 {
			row := make([]int8, m)
			row[i] = 1
			addEq(row, int8(cVal[i]-1))
		}
	}
	// triangle constraints: sum of three edges == 0 mod 3
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			for k := 0; k < j; k++ {
				a := adj[i][j]
				b := adj[j][k]
				c := adj[k][i]
				if a != -1 && b != -1 && c != -1 {
					row := make([]int8, m)
					row[a] = 1
					row[b] = 1
					row[c] = 1
					addEq(row, 0)
				}
			}
		}
	}
	// Gaussian elimination mod 3
	nEq := len(eqA)
	ind := 0
	mod3 := func(x int8) int8 {
		x %= 3
		if x < 0 {
			x += 3
		}
		return x
	}
	for j := 0; j < m && ind < nEq; j++ {
		// find pivot
		piv := -1
		for i := ind; i < nEq; i++ {
			if eqA[i][j] != 0 {
				piv = i
				break
			}
		}
		if piv < 0 {
			continue
		}
		// swap to row ind
		eqA[ind], eqA[piv] = eqA[piv], eqA[ind]
		values[ind], values[piv] = values[piv], values[ind]
		// normalize pivot to 1
		k := invThree(eqA[ind][j])
		if k != 1 {
			for jj := 0; jj < m; jj++ {
				eqA[ind][jj] = mod3(eqA[ind][jj] * k)
			}
			values[ind] = mod3(values[ind] * k)
		}
		// eliminate other rows
		for i := 0; i < nEq; i++ {
			if i == ind {
				continue
			}
			v := eqA[i][j]
			if v != 0 {
				kk := mod3(3 - v)
				for jj := 0; jj < m; jj++ {
					eqA[i][jj] = mod3(eqA[i][jj] + eqA[ind][jj]*kk)
				}
				values[i] = mod3(values[i] + values[ind]*kk)
			}
		}
		ind++
	}
	// extract solution
	ans := make([]int8, m)
	for j := 0; j < m; j++ {
		for i := 0; i < nEq; i++ {
			if eqA[i][j] != 0 {
				ans[j] = values[i]
				// mark used
				values[i] = 0
				break
			}
		}
	}
	// check consistency
	for i := 0; i < len(values); i++ {
		if values[i] != 0 {
			fmt.Fprintln(wrtr, -1)
			return
		}
	}
	// output
	for j := 0; j < m; j++ {
		fmt.Fprint(wrtr, int(ans[j])+1)
		if j+1 < m {
			fmt.Fprint(wrtr, " ")
		}
	}
	fmt.Fprintln(wrtr)
}

func main() {
	defer wrtr.Flush()
	var t int
	fmt.Fscan(rdr, &t)
	for i := 0; i < t; i++ {
		solveOne()
	}
}
