package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		a := make([]int, n)
		orAll := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			orAll |= a[i]
		}
		// Precompute next occurrence for each bit on doubled array
		B := 30
		a2 := make([]int, 2*n)
		copy(a2, a)
		copy(a2[n:], a)
		nxt := make([][]int, B)
		for b := 0; b < B; b++ {
			nxt[b] = make([]int, 2*n+1)
			nxt[b][2*n] = 2 * n
		}
		for i := 2*n - 1; i >= 0; i-- {
			val := a2[i]
			for b := 0; b < B; b++ {
				nxt[b][i] = nxt[b][i+1]
			}
			for b := 0; b < B; b++ {
				if (val>>b)&1 == 1 {
					nxt[b][i] = i
				}
			}
		}
		// Map from OR value to minimal index
		mp := make(map[int]int)
		for s := 0; s < n; s++ {
			pos := s
			length := 1
			orv := a[pos]
			idx := (length-1)*n + (s + 1)
			if cur, ok := mp[orv]; !ok || idx < cur {
				mp[orv] = idx
			}
			mask := orAll &^ orv
			for mask != 0 && length < n {
				nextpos := 2 * n
				mm := mask
				for mm != 0 {
					lb := mm & -mm
					bit := bits.TrailingZeros(uint(lb))
					np := nxt[bit][pos+1]
					if np < nextpos {
						nextpos = np
					}
					mm -= lb
				}
				if nextpos >= s+n {
					break
				}
				pos = nextpos
				length = pos - s + 1
				orv |= a2[pos]
				idx = (length-1)*n + (s + 1)
				if cur, ok := mp[orv]; !ok || idx < cur {
					mp[orv] = idx
				}
				mask = orAll &^ orv
			}
			if orv != orAll {
				idx = (n-1)*n + (s + 1)
				if cur, ok := mp[orAll]; !ok || idx < cur {
					mp[orAll] = idx
				}
			}
		}
		vals := make([]int, 0, len(mp))
		for v := range mp {
			vals = append(vals, v)
		}
		sort.Ints(vals)
		suf := make([]int, len(vals))
		for i := len(vals) - 1; i >= 0; i-- {
			val := mp[vals[i]]
			if i == len(vals)-1 || suf[i+1] > val {
				suf[i] = val
			} else {
				suf[i] = suf[i+1]
			}
		}
		for ; q > 0; q-- {
			var v int
			fmt.Fscan(in, &v)
			j := sort.Search(len(vals), func(i int) bool { return vals[i] > v })
			if j == len(vals) {
				fmt.Fprintln(out, -1)
			} else {
				fmt.Fprintln(out, suf[j])
			}
		}
	}
}
