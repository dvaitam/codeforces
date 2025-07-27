package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const INF int64 = 1 << 60

// one-dimensional squared Euclidean distance transform
func distTransform1D(f []int64) []int64 {
	n := len(f)
	d := make([]int64, n)
	v := make([]int, n)
	z := make([]float64, n+1)
	k := 0
	v[0] = 0
	z[0] = math.Inf(-1)
	z[1] = math.Inf(1)
	for q := 1; q < n; q++ {
		s := ((float64(f[q]) + float64(q*q)) - (float64(f[v[k]]) + float64(v[k]*v[k]))) / (2.0 * float64(q-v[k]))
		for s <= z[k] {
			k--
			s = ((float64(f[q]) + float64(q*q)) - (float64(f[v[k]]) + float64(v[k]*v[k]))) / (2.0 * float64(q-v[k]))
		}
		k++
		v[k] = q
		z[k] = s
		z[k+1] = math.Inf(1)
	}
	k = 0
	for q := 0; q < n; q++ {
		for z[k+1] < float64(q) {
			k++
		}
		diff := q - v[k]
		d[q] = int64(diff*diff) + f[v[k]]
	}
	return d
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	rows := n + 1
	cols := m + 1
	grid := make([][]int64, rows)
	for i := 0; i < rows; i++ {
		grid[i] = make([]int64, cols)
		for j := 0; j < cols; j++ {
			var x int
			fmt.Fscan(in, &x)
			if x == 1 {
				grid[i][j] = 0
			} else {
				grid[i][j] = INF
			}
		}
	}
	// horizontal pass
	horiz := make([][]int64, rows)
	for i := 0; i < rows; i++ {
		horiz[i] = distTransform1D(grid[i])
	}
	// vertical pass
	dist := make([][]int64, rows)
	for j := 0; j < cols; j++ {
		col := make([]int64, rows)
		for i := 0; i < rows; i++ {
			col[i] = horiz[i][j]
		}
		colRes := distTransform1D(col)
		for i := 0; i < rows; i++ {
			if j == 0 {
				dist[i] = make([]int64, cols)
			}
			dist[i][j] = colRes[i]
		}
	}
	var sum int64
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			sum += dist[i][j]
		}
	}
	fmt.Println(sum)
}
