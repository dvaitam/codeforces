package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const maxE = 500005

var (
	a   [maxE]int
	g   [maxE][]int
	S   [maxE]int64
	rdr = bufio.NewReader(os.Stdin)
	wtr = bufio.NewWriter(os.Stdout)
)

func add(l, r int, x int64) {
	if l > r {
		return
	}
	S[l] += x
	if r+1 < maxE {
		S[r+1] -= x
	}
}

type orderedSet struct{ arr []int }

func (s *orderedSet) Insert(x int) int {
	i := sort.SearchInts(s.arr, x)
	if i == len(s.arr) || s.arr[i] != x {
		s.arr = append(s.arr, 0)
		copy(s.arr[i+1:], s.arr[i:])
		s.arr[i] = x
	}
	return i
}

func (s *orderedSet) At(i int) int { return s.arr[i] }
func (s *orderedSet) Size() int    { return len(s.arr) }

func main() {
	defer wtr.Flush()
	var T int
	fmt.Fscan(rdr, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(rdr, &n)
		S[0] = 0
		for i := 1; i <= n; i++ {
			g[i] = g[i][:0]
			S[i] = 0
		}
		for i := 1; i <= n; i++ {
			fmt.Fscan(rdr, &a[i])
			g[a[i]] = append(g[a[i]], i)
		}
		se := orderedSet{}
		se.Insert(0)
		se.Insert(n + 1)
		for i := 1; i <= n; i++ {
			for _, x := range g[i] {
				idx := se.Insert(x)
				t1, t2 := idx-1, idx+1
				L := se.At(t1)
				R := se.At(t2)
				add(0, L-1, int64(x-L)*int64(R-x)*int64(a[x]))
				add(R+1, n, int64(x-L)*int64(R-x)*int64(a[x]))
				add(L+1, x-1, int64(x-L-1)*int64(R-x)*int64(a[x]))
				add(x+1, R-1, int64(x-L)*int64(R-x-1)*int64(a[x]))
				if L != 0 && t1-1 >= 0 {
					LL := se.At(t1 - 1)
					add(L, L, int64(x-LL-1)*int64(R-x)*int64(a[x]))
				}
				if R != n+1 && t2+1 < se.Size() {
					RR := se.At(t2 + 1)
					add(R, R, int64(x-L)*int64(RR-x-1)*int64(a[x]))
				}
			}
		}
		for i := 1; i <= n; i++ {
			S[i] += S[i-1]
			fmt.Fprint(wtr, S[i], " ")
		}
		fmt.Fprintln(wtr)
	}
}
