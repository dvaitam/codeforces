package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for Codeforces problem 1881G - as described in problemG.txt.
// It supports range additions on a string (modulo 26) and checks whether
// a substring is "beautiful" (contains no palindromic substring of length
// at least two). Only palindromes of length two or three need to be checked.

// Fenwick tree implementation for prefix sums.
type Fenwick struct {
	n   int
	bit []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, bit: make([]int, n+2)}
}

func (f *Fenwick) Add(idx, val int) {
	for idx <= f.n {
		f.bit[idx] += val
		idx += idx & -idx
	}
}

func (f *Fenwick) Prefix(idx int) int {
	sum := 0
	for idx > 0 {
		sum += f.bit[idx]
		idx &= idx - 1
	}
	return sum
}

func (f *Fenwick) RangeSum(l, r int) int {
	if r < l {
		return 0
	}
	return f.Prefix(r) - f.Prefix(l-1)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		var s string
		fmt.Fscan(reader, &s)

		base := make([]int, n+1) // 1-indexed values of characters
		for i := 1; i <= n; i++ {
			base[i] = int(s[i-1] - 'a')
		}

		diff := NewFenwick(n + 2) // range add, point query
		eq1 := make([]int, n)     // for i in [1..n-1]
		eq2 := make([]int, n)     // for i in [1..n-2]
		bit1 := NewFenwick(n - 1)
		bit2 := NewFenwick(n - 2)

		charAt := func(idx int) int {
			val := base[idx] + diff.Prefix(idx)
			val %= 26
			if val < 0 {
				val += 26
			}
			return val
		}

		// Initialize equality arrays and BITs
		for i := 1; i < n; i++ {
			if base[i] == base[i+1] {
				eq1[i] = 1
				bit1.Add(i, 1)
			}
		}
		for i := 1; i+2 <= n; i++ {
			if base[i] == base[i+2] {
				eq2[i] = 1
				bit2.Add(i, 1)
			}
		}

		updateEq1 := func(pos int) {
			if pos < 1 || pos >= n {
				return
			}
			newVal := 0
			if charAt(pos) == charAt(pos+1) {
				newVal = 1
			}
			if newVal != eq1[pos] {
				bit1.Add(pos, newVal-eq1[pos])
				eq1[pos] = newVal
			}
		}

		updateEq2 := func(pos int) {
			if pos < 1 || pos+2 > n {
				return
			}
			newVal := 0
			if charAt(pos) == charAt(pos+2) {
				newVal = 1
			}
			if newVal != eq2[pos] {
				bit2.Add(pos, newVal-eq2[pos])
				eq2[pos] = newVal
			}
		}

		for ; m > 0; m-- {
			var typ int
			fmt.Fscan(reader, &typ)
			if typ == 1 {
				var l, r, x int
				fmt.Fscan(reader, &l, &r, &x)
				x %= 26
				diff.Add(l, x)
				diff.Add(r+1, -x)
				updateEq1(l - 1)
				updateEq1(r)
				updateEq2(l - 2)
				updateEq2(l - 1)
				updateEq2(r - 1)
				updateEq2(r)
			} else {
				var l, r int
				fmt.Fscan(reader, &l, &r)
				if r-l+1 < 2 {
					fmt.Fprintln(writer, "YES")
					continue
				}
				if bit1.RangeSum(l, r-1) > 0 {
					fmt.Fprintln(writer, "NO")
					continue
				}
				if r-l+1 >= 3 && bit2.RangeSum(l, r-2) > 0 {
					fmt.Fprintln(writer, "NO")
				} else {
					fmt.Fprintln(writer, "YES")
				}
			}
		}
	}
}
