package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// ---------- Fenwick (BIT) ----------
type BIT []int

func (b BIT) add(idx, delta int) {
	for idx < len(b) {
		b[idx] += delta
		idx += idx & -idx
	}
}
func (b BIT) sum(idx int) int {
	res := 0
	for idx > 0 {
		res += b[idx]
		idx &= idx - 1
	}
	return res
}
// first index such that prefix > k (1‑based, always exists if k < total)
func (b BIT) kth(k int) int {
	idx := 0
	bit := 1
	for bit<<1 < len(b) {
		bit <<= 1
	}
	for bit > 0 {
		nxt := idx + bit
		if nxt < len(b) && b[nxt] <= k {
			k -= b[nxt]
			idx = nxt
		}
		bit >>= 1
	}
	return idx + 1
}

// -----------------------------------

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		l := make([]int, n)
		r := make([]int, n)
		uniq := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &l[i], &r[i])
			uniq[i] = l[i]
		}
		sort.Ints(uniq)
		uniq = unique(uniq)

		pos := make(map[int]int, len(uniq))
		for i, v := range uniq {
			pos[v] = i + 1 // 1‑based for BIT
		}

		bit := make(BIT, len(uniq)+2)
		total := 0

		for i := 0; i < n; i++ {
			// 1. remove smallest element > r[i] (if any)
			u := sort.Search(len(uniq), func(j int) bool { return uniq[j] > r[i] }) + 1 // 1‑based
			left := bit.sum(u - 1)
			if total-left > 0 { // there is at least one element > r[i]
				idx := bit.kth(left) // first such position
				bit.add(idx, -1)
				total--
			}

			// 2. insert l[i]
			bit.add(pos[l[i]], 1)
			total++

			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, total)
		}
		fmt.Fprintln(out)
	}
}

func unique(a []int) []int {
	j := 0
	for i := 1; i < len(a); i++ {
		if a[i] != a[j] {
			j++
			a[j] = a[i]
		}
	}
	return a[:j+1]
}

