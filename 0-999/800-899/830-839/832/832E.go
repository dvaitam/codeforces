package main

import (
	"bufio"
	"fmt"
	"os"
)

var inv = [5]int{0, 1, 3, 2, 4}

func insert(basis [][]int, v []int) bool {
	m := len(v)
	for j := 0; j < m; j++ {
		if v[j] == 0 {
			continue
		}
		if basis[j] == nil {
			invv := inv[v[j]]
			for k := j; k < m; k++ {
				v[k] = v[k] * invv % 5
			}
			b := make([]int, m)
			copy(b, v)
			basis[j] = b
			return true
		} else {
			f := v[j]
			for k := j; k < m; k++ {
				v[k] -= f * basis[j][k]
				v[k] %= 5
				if v[k] < 0 {
					v[k] += 5
				}
			}
		}
	}
	return false
}
func inSpan(basis [][]int, v []int) bool {
	m := len(v)
	for j := 0; j < m; j++ {
		if v[j] == 0 {
			continue
		}
		if basis[j] == nil {
			return false
		}
		f := v[j]
		for k := j; k < m; k++ {
			v[k] -= f * basis[j][k]
			v[k] %= 5
			if v[k] < 0 {
				v[k] += 5
			}
		}
	}
	return true
}
func powmod(a, b int64) int64 {
	mod := int64(1e9 + 7)
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}
func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m)
	basis := make([][]int, m)
	rank := 0
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		v := make([]int, m)
		for j := 0; j < m; j++ {
			v[j] = int(s[j]-'a') % 5
		}
		if insert(basis, v) {
			rank++
		}
	}
	pow := powmod(5, int64(n-rank))
	var q int
	fmt.Fscan(in, &q)
	out := bufio.NewWriter(os.Stdout)
	for ; q > 0; q-- {
		var b string
		fmt.Fscan(in, &b)
		v := make([]int, m)
		for j := 0; j < m; j++ {
			v[j] = int(b[j] - 'a')
		}
		if inSpan(basis, v) {
			fmt.Fprintln(out, pow)
		} else {
			fmt.Fprintln(out, 0)
		}
	}
	out.Flush()
}
