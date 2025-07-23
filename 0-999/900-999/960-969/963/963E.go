package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func modInv(a int64) int64 {
	return modPow(a%MOD, MOD-2)
}

// solve linear system using Gaussian elimination modulo MOD
func gauss(a [][]int64, n int) []int64 {
	for i := 0; i < n; i++ {
		piv := i
		for piv < n && a[piv][i] == 0 {
			piv++
		}
		if piv == n {
			continue
		}
		if piv != i {
			a[i], a[piv] = a[piv], a[i]
		}
		inv := modInv(a[i][i])
		for j := i; j <= n; j++ {
			a[i][j] = a[i][j] * inv % MOD
		}
		for k := 0; k < n; k++ {
			if k == i || a[k][i] == 0 {
				continue
			}
			factor := a[k][i]
			for j := i; j <= n; j++ {
				a[k][j] = (a[k][j] - factor*a[i][j]) % MOD
				if a[k][j] < 0 {
					a[k][j] += MOD
				}
			}
		}
	}
	ans := make([]int64, n)
	for i := 0; i < n; i++ {
		ans[i] = a[i][n] % MOD
		if ans[i] < 0 {
			ans[i] += MOD
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var R int
	var a1, a2, a3, a4 int64
	if _, err := fmt.Fscan(in, &R, &a1, &a2, &a3, &a4); err != nil {
		return
	}
	sum := a1 + a2 + a3 + a4
	invSum := modInv(sum % MOD)
	p := []int64{a1 % MOD * invSum % MOD, a2 % MOD * invSum % MOD, a3 % MOD * invSum % MOD, a4 % MOD * invSum % MOD}

	type pt struct{ x, y int }
	points := make([]pt, 0)
	idx := make(map[pt]int)
	for x := -R; x <= R; x++ {
		for y := -R; y <= R; y++ {
			if x*x+y*y <= R*R {
				idx[pt{x, y}] = len(points)
				points = append(points, pt{x, y})
			}
		}
	}
	n := len(points)
	mat := make([][]int64, n)
	for i := 0; i < n; i++ {
		row := make([]int64, n+1)
		row[i] = 1
		x := points[i].x
		y := points[i].y
		dirs := []pt{{-1, 0}, {0, -1}, {1, 0}, {0, 1}}
		for d, dv := range dirs {
			nx, ny := x+dv.x, y+dv.y
			if j, ok := idx[pt{nx, ny}]; ok {
				row[j] = (row[j] - p[d]) % MOD
			}
		}
		row[n] = 1
		mat[i] = row
	}
	sol := gauss(mat, n)
	if id, ok := idx[pt{0, 0}]; ok {
		fmt.Println(sol[id] % MOD)
	}
}
