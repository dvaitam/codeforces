package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD1 int64 = 1000000007
const MOD2 int64 = 1000000009
const BASE int64 = 911382323

var (
	n       int
	g       [][]int
	dpDown1 []int64
	dpDown2 []int64
	dpUp1   []int64
	dpUp2   []int64
)

func dfsDown(v, p int) {
	dpDown1[v] = 1
	dpDown2[v] = 1
	for _, to := range g[v] {
		if to == p {
			continue
		}
		dfsDown(to, v)
		dpDown1[v] = (dpDown1[v] + BASE*dpDown1[to]%MOD1) % MOD1
		dpDown2[v] = (dpDown2[v] + BASE*dpDown2[to]%MOD2) % MOD2
	}
}

func dfsUp(v, p int) {
	for _, to := range g[v] {
		if to == p {
			continue
		}
		dpUp1[to] = BASE * (dpUp1[v] + dpDown1[v] - BASE*dpDown1[to]%MOD1 + MOD1) % MOD1
		dpUp2[to] = BASE * (dpUp2[v] + dpDown2[v] - BASE*dpDown2[to]%MOD2 + MOD2) % MOD2
		dfsUp(to, v)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n)
	a := make([]int, n-1)
	count0 := 0
	for i := 0; i < n-1; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] == 0 {
			count0++
		}
	}
	if count0 > 1 {
		fmt.Fprintln(writer, 0)
		fmt.Fprintln(writer)
		return
	}
	g = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	pow1 := make([]int64, n)
	pow2 := make([]int64, n)
	pow1[0] = 1
	pow2[0] = 1
	for i := 1; i < n; i++ {
		pow1[i] = pow1[i-1] * BASE % MOD1
		pow2[i] = pow2[i-1] * BASE % MOD2
	}
	HA1, HA2 := int64(0), int64(0)
	for _, x := range a {
		HA1 = (HA1 + pow1[x]) % MOD1
		HA2 = (HA2 + pow2[x]) % MOD2
	}

	dpDown1 = make([]int64, n+1)
	dpDown2 = make([]int64, n+1)
	dpUp1 = make([]int64, n+1)
	dpUp2 = make([]int64, n+1)
	dfsDown(1, 0)
	dfsUp(1, 0)

	H1 := make([]int64, n+1)
	H2 := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		H1[i] = (dpDown1[i] + dpUp1[i]) % MOD1
		H2[i] = (dpDown2[i] + dpUp2[i]) % MOD2
	}

	powIndex := make(map[int64]int)
	for d := 0; d < n; d++ {
		powIndex[pow1[d]] = d
	}

	good := make([]int, 0)
	for i := 1; i <= n; i++ {
		diff1 := H1[i] - HA1
		if diff1 < 0 {
			diff1 += MOD1
		}
		d, ok := powIndex[diff1]
		if !ok {
			continue
		}
		diff2 := H2[i] - HA2
		if diff2 < 0 {
			diff2 += MOD2
		}
		if diff2 != pow2[d] {
			continue
		}
		if (count0 == 0 && d == 0) || (count0 == 1 && d != 0) {
			good = append(good, i)
		}
	}

	fmt.Fprintln(writer, len(good))
	for i, idx := range good {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, idx)
	}
	if len(good) > 0 {
		writer.WriteByte('\n')
	}
}
