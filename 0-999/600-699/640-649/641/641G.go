package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 1000000007

func powMod(a, b int) int {
	res := 1
	base := int64(a)
	for b > 0 {
		if b&1 == 1 {
			res = int((int64(res) * base) % MOD)
		}
		base = (base * base) % MOD
		b >>= 1
	}
	return res
}

func modInverse(n int) int {
	return powMod(n, MOD-2)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	scanInt := func() int {
		scanner.Scan()
		s := scanner.Text()
		x := 0
		for i := 0; i < len(s); i++ {
			x = x*10 + int(s[i]-'0')
		}
		return x
	}

	if !scanner.Scan() {
		return
	}
	nStr := scanner.Text()
	n := 0
	for i := 0; i < len(nStr); i++ {
		n = n*10 + int(nStr[i]-'0')
	}

	k := scanInt()

	diag := make([]int, n+1)
	mat := make([]map[int]int, n+1)
	for i := 1; i <= n; i++ {
		mat[i] = make(map[int]int)
	}

	// Initialize initial clique (vertices 1 to k)
	for i := 1; i <= k; i++ {
		diag[i] = k - 1
		for j := 1; j < i; j++ {
			mat[i][j] = MOD - 1
		}
	}

	// Read vertices k+1 to n
	for i := k + 1; i <= n; i++ {
		diag[i] = k
		for j := 0; j < k; j++ {
			v := scanInt()
			// Edge between i and v (v < i)
			mat[i][v] = MOD - 1
			diag[v]++
		}
	}

	ans := int64(1)

	// Eliminate vertices from n down to 2
	for i := n; i >= 2; i-- {
		d := diag[i] % MOD
		if d < 0 {
			d += MOD
		}
		ans = (ans * int64(d)) % MOD

		inv := modInverse(d)

		neighbors := make([]int, 0, len(mat[i]))
		for u := range mat[i] {
			neighbors = append(neighbors, u)
		}

		for idx, u := range neighbors {
			valU := mat[i][u]

			// Update diag[u]
			term := int((int64(valU) * int64(valU)) % MOD)
			term = int((int64(term) * int64(inv)) % MOD)
			diag[u] = (diag[u] - term + MOD) % MOD

			// Update edges between neighbors
			for j := 0; j < idx; j++ {
				v := neighbors[j]

				var big, small int
				if u > v {
					big, small = u, v
				} else {
					big, small = v, u
				}

				valV := mat[i][v]

				update := int((int64(valU) * int64(valV)) % MOD)
				update = int((int64(update) * int64(inv)) % MOD)

				old := mat[big][small]
				mat[big][small] = (old - update + MOD) % MOD
			}
		}
	}

	fmt.Println(ans)
}
