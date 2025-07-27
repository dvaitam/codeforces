package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 31607

func modPow(a, b int) int {
	res := 1
	a %= MOD
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func modInv(a int) int {
	return modPow(a, MOD-2)
}

func fwtAnd(a []int, inv bool) {
	n := len(a)
	for step := 1; step < n; step <<= 1 {
		for i := 0; i < n; i += step << 1 {
			for j := 0; j < step; j++ {
				if inv {
					a[i+j] = (a[i+j] - a[i+j+step] + MOD) % MOD
				} else {
					a[i+j] = (a[i+j] + a[i+j+step]) % MOD
				}
			}
		}
	}
}

// probability that there is no fully successful row or column
func noRowNoCol(p [][]int) int {
	n := len(p)
	size := 1 << uint(n)
	rows := make([][]int, n)
	for i := 0; i < n; i++ {
		arr := make([]int, size)
		arr[0] = 1
		for j := 0; j < n; j++ {
			next := make([]int, size)
			pj := p[i][j]
			qj := (1 - pj + MOD) % MOD
			for mask := 0; mask < size; mask++ {
				if arr[mask] == 0 {
					continue
				}
				next[mask] = (next[mask] + arr[mask]*qj) % MOD
				next[mask|1<<uint(j)] = (next[mask|1<<uint(j)] + arr[mask]*pj) % MOD
			}
			arr = next
		}
		// row completely successful leads to winning table -> exclude
		arr[size-1] = 0
		fwtAnd(arr, false)
		rows[i] = arr
	}
	res := make([]int, size)
	for mask := 0; mask < size; mask++ {
		res[mask] = 1
	}
	for i := 0; i < n; i++ {
		for mask := 0; mask < size; mask++ {
			res[mask] = res[mask] * rows[i][mask] % MOD
		}
	}
	fwtAnd(res, true)
	// mask==0 means every column has at least one failure
	return res[0]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	inv10000 := modInv(10000)
	orig := make([][]int, n)
	for i := 0; i < n; i++ {
		orig[i] = make([]int, n)
		for j := 0; j < n; j++ {
			var a int
			fmt.Fscan(in, &a)
			orig[i][j] = a * inv10000 % MOD
		}
	}

	type diagSel struct{ main, anti bool }
	selections := []diagSel{{false, false}, {true, false}, {false, true}, {true, true}}
	results := make([]int, 4)

	for idx, s := range selections {
		// copy matrix and apply forced successes
		mat := make([][]int, n)
		for i := 0; i < n; i++ {
			mat[i] = make([]int, n)
			copy(mat[i], orig[i])
		}
		factor := 1
		for i := 0; i < n; i++ {
			if s.main {
				factor = factor * mat[i][i] % MOD
				mat[i][i] = 1
			}
			if s.anti {
				j := n - 1 - i
				if !(s.main && j == i) {
					factor = factor * mat[i][j] % MOD
					mat[i][j] = 1
				}
			}
		}
		val := noRowNoCol(mat)
		results[idx] = factor * val % MOD
	}

	probNoLine := (results[0] - results[1] - results[2] + results[3]) % MOD
	if probNoLine < 0 {
		probNoLine += MOD
	}
	ans := (1 - probNoLine) % MOD
	if ans < 0 {
		ans += MOD
	}
	fmt.Println(ans)
}
