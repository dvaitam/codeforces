package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int = 1e9 + 7

// Basis represents a linear basis for xor operations up to 20 bits.
type Basis struct {
	b    [20]int
	size int
}

func (bs *Basis) Add(x int) {
	for i := 19; i >= 0; i-- {
		if (x>>i)&1 == 0 {
			continue
		}
		if bs.b[i] != 0 {
			x ^= bs.b[i]
		} else {
			bs.b[i] = x
			bs.size++
			return
		}
	}
}

func (bs *Basis) Contains(x int) bool {
	for i := 19; i >= 0; i-- {
		if (x>>i)&1 == 0 {
			continue
		}
		if bs.b[i] == 0 {
			return false
		}
		x ^= bs.b[i]
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	type Query struct {
		x   int
		idx int
	}
	queries := make([][]Query, n+1)
	for i := 0; i < q; i++ {
		var l, x int
		fmt.Fscan(reader, &l, &x)
		queries[l] = append(queries[l], Query{x, i})
	}

	pow2 := make([]int, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = pow2[i-1] * 2 % mod
	}

	ans := make([]int, q)
	var bs Basis
	for i := 1; i <= n; i++ {
		bs.Add(arr[i])
		for _, qu := range queries[i] {
			if bs.Contains(qu.x) {
				ans[qu.idx] = pow2[i-bs.size]
			} else {
				ans[qu.idx] = 0
			}
		}
	}
	for i := 0; i < q; i++ {
		fmt.Fprintln(writer, ans[i])
	}
}
