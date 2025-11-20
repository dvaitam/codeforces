package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const MOD = 1000000007

func power(a, b int64) int64 {
	var res int64 = 1
	a %= MOD
	for b > 0 {
		if b&1 == 1 {
			res = (res * a) % MOD
		}
		a = (a * a) % MOD
		b >>= 1
	}
	return res
}

func modInverse(n int64) int64 {
	return power(n, MOD-2)
}

func gaussianDet(a [][]int64, n int) int64 {
	var res int64 = 1
	mat := make([][]int64, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]int64, n)
		copy(mat[i], a[i])
	}

	for i := 0; i < n; i++ {
		pivot := i
		for j := i + 1; j < n; j++ {
			if mat[j][i] != 0 {
				pivot = j
				break
			}
		}
		if mat[pivot][i] == 0 {
			return 0
		}
		if pivot != i {
			mat[i], mat[pivot] = mat[pivot], mat[i]
			res = (MOD - res) % MOD
		}
		res = (res * mat[i][i]) % MOD
		inv := modInverse(mat[i][i])

		for j := i + 1; j < n; j++ {
			if mat[j][i] != 0 {
				factor := (mat[j][i] * inv) % MOD
				for l := i; l < n; l++ {
					sub := (factor * mat[i][l]) % MOD
					mat[j][l] = (mat[j][l] - sub + MOD) % MOD
				}
			}
		}
	}
	return res
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	scanInt := func() int {
		scanner.Scan()
		val, _ := strconv.Atoi(scanner.Text())
		return val
	}

	n := scanInt()
	k := scanInt()

	neighbors := make([][]int, n+1)
	for i := k + 1; i <= n; i++ {
		neighbors[i] = make([]int, k)
		for j := 0; j < k; j++ {
			neighbors[i][j] = scanInt()
		}
	}

	adj := make([]map[int]int64, n+1)
	for i := 1; i <= n; i++ {
		adj[i] = make(map[int]int64)
	}

	for i := 1; i <= k; i++ {
		for j := i + 1; j <= k; j++ {
			adj[i][j] = 1
			adj[j][i] = 1
		}
	}
	for i := k + 1; i <= n; i++ {
		for _, u := range neighbors[i] {
			adj[i][u] = 1
			adj[u][i] = 1
		}
	}

	var ans int64 = 1

	for i := n; i > k; i-- {
		var D int64 = 0
		for _, u := range neighbors[i] {
			D = (D + adj[i][u]) % MOD
		}

		ans = (ans * D) % MOD
		invD := modInverse(D)

		ns := neighbors[i]
		for a := 0; a < k; a++ {
			u := ns[a]
			wu := adj[i][u]
			for b := a + 1; b < k; b++ {
				v := ns[b]
				wv := adj[i][v]

				term := (wu * wv) % MOD
				term = (term * invD) % MOD

				old := adj[u][v]
				newVal := (old + term) % MOD
				adj[u][v] = newVal
				adj[v][u] = newVal
			}
		}

		adj[i] = nil
		for _, u := range neighbors[i] {
			delete(adj[u], i)
		}
	}

	if k == 1 {
		fmt.Println(ans)
		return
	}

	mat := make([][]int64, k)
	for i := range mat {
		mat[i] = make([]int64, k)
	}

	for i := 1; i <= k; i++ {
		var sum int64 = 0
		for j := 1; j <= k; j++ {
			if i == j {
				continue
			}
			w := adj[i][j]
			mat[i-1][j-1] = (MOD - w) % MOD
			sum = (sum + w) % MOD
		}
		mat[i-1][i-1] = sum
	}

	det := gaussianDet(mat, k-1)
	ans = (ans * det) % MOD

	fmt.Println(ans)
}
