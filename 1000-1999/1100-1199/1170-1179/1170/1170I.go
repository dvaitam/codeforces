package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const mod = 998244353

type SegTree struct {
   n    int
   sum  []int
   lazy []int
}

func newSegTree(n int) *SegTree {
   size := 1
   for size < n*4 {
       size <<= 1
   }
   st := &SegTree{
       n:    n,
       sum:  make([]int, size),
       lazy: make([]int, size),
   }
   // init lazy to 1
   for i := range st.lazy {
       st.lazy[i] = 1
   }
   return st
}

func (st *SegTree) applyMul(v, mul int) {
   st.sum[v] = int(int64(st.sum[v]) * int64(mul) % mod)
   st.lazy[v] = int(int64(st.lazy[v]) * int64(mul) % mod)
}

func (st *SegTree) push(v int) {
   if st.lazy[v] != 1 {
       st.applyMul(v*2, st.lazy[v])
       st.applyMul(v*2+1, st.lazy[v])
       st.lazy[v] = 1
   }
}

func (st *SegTree) pull(v int) {
   st.sum[v] = st.sum[v*2] + st.sum[v*2+1]
   if st.sum[v] >= mod {
       st.sum[v] -= mod
   }
}

// range multiply [l..r] by mul
func (st *SegTree) rangeMul(v, tl, tr, l, r, mul int) {
   if l > r {
       return
   }
   if l <= tl && tr <= r {
       st.applyMul(v, mul)
       return
   }
   st.push(v)
   tm := (tl + tr) >> 1
   if l <= tm {
       st.rangeMul(v*2, tl, tm, l, r, mul)
   }
   if r > tm {
       st.rangeMul(v*2+1, tm+1, tr, l, r, mul)
   }
   st.pull(v)
}

// sum query [l..r]
func (st *SegTree) rangeSum(v, tl, tr, l, r int) int {
   if l > r {
       return 0
   }
   if l <= tl && tr <= r {
       return st.sum[v]
   }
   st.push(v)
   tm := (tl + tr) >> 1
   res := 0
   if l <= tm {
       res = st.rangeSum(v*2, tl, tm, l, r)
   }
   if r > tm {
       res2 := st.rangeSum(v*2+1, tm+1, tr, l, r)
       res += res2
       if res >= mod {
           res -= mod
       }
   }
   return res
}

// point add at pos idx: add val
func (st *SegTree) pointAdd(v, tl, tr, idx, val int) {
   if tl == tr {
       st.sum[v] += val
       if st.sum[v] >= mod {
           st.sum[v] -= mod
       }
       return
   }
   st.push(v)
   tm := (tl + tr) >> 1
   if idx <= tm {
       st.pointAdd(v*2, tl, tm, idx, val)
   } else {
       st.pointAdd(v*2+1, tm+1, tr, idx, val)
   }
   st.pull(v)
}


func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   segs := make([][2]int, n)
   coords := make([]int, 0, n*2)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &segs[i][0], &segs[i][1])
       coords = append(coords, segs[i][0], segs[i][1])
   }
   sort.Ints(coords)
   coords = unique(coords)
   mcoord := len(coords)
   // map to coord indices
   Lold := make([]int, n)
   Rold := make([]int, n)
   for i := 0; i < n; i++ {
       l, r := segs[i][0], segs[i][1]
       Lold[i] = sort.SearchInts(coords, l)
       Rold[i] = sort.SearchInts(coords, r)
   }
   // mark coverage of interval cells and point cells
   // points at coords indices 0..mcoord-1, intervals at 0..mcoord-2
   coverPoint := make([]int, mcoord+1)
   coverInterval := make([]int, mcoord+1)
   for i := 0; i < n; i++ {
       // point coverage
       coverPoint[Lold[i]]++
       coverPoint[Rold[i]+1]--
       // interval coverage
       if Lold[i] < Rold[i] {
           coverInterval[Lold[i]]++
           coverInterval[Rold[i]]--
       }
   }
   // compute valid point and interval cells
   pointCov := make([]bool, mcoord)
   intervalCov := make([]bool, mcoord)
   curP, curI := 0, 0
   for j := 0; j < mcoord; j++ {
       curP += coverPoint[j]
       pointCov[j] = curP > 0
       if j < mcoord-1 {
           curI += coverInterval[j]
           intervalCov[j] = curI > 0
       }
   }
   // build unified cells: even c=2*j point, odd c=2*j+1 interval
   cells := 2*mcoord - 1
   validCell := make([]bool, cells)
   for c := 0; c < cells; c++ {
       if c%2 == 0 {
           j := c / 2
           validCell[c] = pointCov[j]
       } else {
           j := (c - 1) / 2
           validCell[c] = intervalCov[j]
       }
   }
   // map old cells to new indices
   oldToNewCell := make([]int, cells)
   for i := range oldToNewCell {
       oldToNewCell[i] = -1
   }
   m := 0
   for c := 0; c < cells; c++ {
       if validCell[c] {
           m++
           oldToNewCell[c] = m
       }
   }
   // build next and prev valid cell lookup
   prevValid := make([]int, cells)
   nextValid := make([]int, cells)
   last := -1
   for c := 0; c < cells; c++ {
       if validCell[c] {
           last = c
       }
       prevValid[c] = last
   }
   next := -1
   for c := cells - 1; c >= 0; c-- {
       if validCell[c] {
           next = c
       }
       nextValid[c] = next
   }
   // prepare DP segments covering cells
   type iv struct{ l, r int }
   ivs := make([]iv, 0, n)
   for i := 0; i < n; i++ {
       left := 2 * Lold[i]
       right := 2 * Rold[i]
       if left < 0 || right >= cells {
           continue
       }
       c1 := nextValid[left]
       c2 := prevValid[right]
       if c1 < 0 || c2 < 0 || c1 > c2 {
           continue
       }
       lnew := oldToNewCell[c1]
       rnew := oldToNewCell[c2]
       if lnew <= 0 || rnew <= 0 {
           continue
       }
       ivs = append(ivs, iv{lnew, rnew})
   }
   // DP
   st := newSegTree(m + 1)
   // dp[0] = 1
   st.pointAdd(1, 0, st.n, 0, 1)
   for _, seg := range ivs {
       l, r := seg.l, seg.r
       // sum dp[l-1..r]
       s := st.rangeSum(1, 0, st.n, l-1, r)
       // multiply [0..l-2] and [r+1..m] by 2
       if l-2 >= 0 {
           st.rangeMul(1, 0, st.n, 0, l-2, 2)
       }
       if r+1 <= m {
           st.rangeMul(1, 0, st.n, r+1, m, 2)
       }
       // add sum to dp[r]
       st.pointAdd(1, 0, st.n, r, s)
   }
   ans := st.rangeSum(1, 0, st.n, m, m)
   fmt.Println(ans)
}

func unique(a []int) []int {
   j := 0
   for i := 0; i < len(a); i++ {
       if i == 0 || a[i] != a[i-1] {
           a[j] = a[i]
           j++
       }
   }
   return a[:j]
}
