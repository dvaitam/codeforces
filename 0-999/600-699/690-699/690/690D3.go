package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000003

func matMul(a, b [][]int64) [][]int64 {
	n := len(a)
	res := make([][]int64, n)
	for i := range res {
		res[i] = make([]int64, n)
	}
	for i := 0; i < n; i++ {
		for k := 0; k < n; k++ {
			if a[i][k] == 0 {
				continue
			}
			av := a[i][k]
			for j := 0; j < n; j++ {
				if b[k][j] == 0 {
					continue
				}
				res[i][j] = (res[i][j] + av*b[k][j]) % MOD
			}
		}
	}
	return res
}

func matPow(mat [][]int64, exp int64) [][]int64 {
	n := len(mat)
	res := make([][]int64, n)
	for i := range res {
		res[i] = make([]int64, n)
		res[i][i] = 1
	}
	for exp > 0 {
		if exp&1 == 1 {
			res = matMul(res, mat)
		}
		mat = matMul(mat, mat)
		exp >>= 1
	}
	return res
}

func vecMul(v []int64, m [][]int64) []int64 {
	n := len(v)
	res := make([]int64, n)
	for j := 0; j < n; j++ {
		var s int64
		for k := 0; k < n; k++ {
			if v[k] == 0 || m[k][j] == 0 {
				continue
			}
			s = (s + v[k]*m[k][j]) % MOD
		}
		res[j] = s
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var C int64
	var W, H int
	if _, err := fmt.Fscan(reader, &C, &W, &H); err != nil {
		return
	}
	n := W + 1
	mat := make([][]int64, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]int64, n)
	}
	for i := 0; i < n; i++ {
		mat[i][0] = 1
		if i < W {
			mat[i][i+1] = int64(H) % MOD
		}
	}
	powMat := matPow(mat, C)
	vec := make([]int64, n)
	vec[0] = 1
	vec = vecMul(vec, powMat)
	var ans int64
	for _, v := range vec {
		ans = (ans + v) % MOD
	}
	fmt.Println(ans)
}
