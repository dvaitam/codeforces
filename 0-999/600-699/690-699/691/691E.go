package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const mod int64 = 1000000007

func matMul(a, b [][]int64) [][]int64 {
	n := len(a)
	res := make([][]int64, n)
	for i := range res {
		res[i] = make([]int64, n)
	}
	for i := 0; i < n; i++ {
		ai := a[i]
		ri := res[i]
		for k := 0; k < n; k++ {
			if ai[k] == 0 {
				continue
			}
			aik := ai[k]
			bk := b[k]
			for j := 0; j < n; j++ {
				if bk[j] == 0 {
					continue
				}
				ri[j] = (ri[j] + aik*bk[j]) % mod
			}
		}
	}
	return res
}

func matVecMul(a [][]int64, v []int64) []int64 {
	n := len(a)
	res := make([]int64, n)
	for i := 0; i < n; i++ {
		row := a[i]
		var sum int64
		for j := 0; j < n; j++ {
			if row[j] == 0 || v[j] == 0 {
				continue
			}
			sum = (sum + row[j]*v[j]) % mod
		}
		res[i] = sum
	}
	return res
}

func powMatVec(m [][]int64, e int64, vec []int64) []int64 {
	for e > 0 {
		if e&1 == 1 {
			vec = matVecMul(m, vec)
		}
		m = matMul(m, m)
		e >>= 1
	}
	return vec
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	if k == 1 {
		fmt.Fprintln(writer, n%int(mod))
		return
	}

	adj := make([][]int64, n)
	for i := 0; i < n; i++ {
		adj[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			if bits.OnesCount64(uint64(arr[i]^arr[j]))%3 == 0 {
				adj[i][j] = 1
			}
		}
	}

	vec := make([]int64, n)
	for i := range vec {
		vec[i] = 1
	}
	vec = powMatVec(adj, k-1, vec)
	var ans int64
	for _, v := range vec {
		ans = (ans + v) % mod
	}
	fmt.Fprintln(writer, ans)
}
