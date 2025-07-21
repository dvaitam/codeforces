package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// BIT is a Fenwick tree for range sum [1..n]
type BIT struct {
   n int
   t []int
}

func NewBIT(n int) *BIT {
   return &BIT{n: n, t: make([]int, n+1)}
}

// Add adds v at index i (1-based)
func (b *BIT) Add(i, v int) {
   for ; i <= b.n; i += i & -i {
       b.t[i] += v
   }
}

// Sum returns sum of [1..i]
func (b *BIT) Sum(i int) int {
   if i > b.n {
       i = b.n
   }
   s := 0
   for ; i > 0; i -= i & -i {
       s += b.t[i]
   }
   return s
}

// Segment holds a segment [l,r]
type Segment struct{ l, r int }

// GapQuery holds a query for segments with l>x and r<y
type GapQuery struct{ x, y, idx int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   segs := make([]Segment, n)
   lArr := make([]int, n)
   rArr := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &segs[i].l, &segs[i].r)
       lArr[i] = segs[i].l
       rArr[i] = segs[i].r
   }
   sort.Ints(lArr)
   sort.Ints(rArr)

   ansL := make([]int, m)
   ansR := make([]int, m)
   gapSum := make([]int, m)
   var gaps []GapQuery

   for qi := 0; qi < m; qi++ {
       var cnt int
       fmt.Fscan(reader, &cnt)
       pts := make([]int, cnt)
       for i := 0; i < cnt; i++ {
           fmt.Fscan(reader, &pts[i])
       }
       // count segments with r < pts[0]
       p1 := pts[0]
       // upper_bound rArr for p1-1
       cL := sort.Search(n, func(i int) bool { return rArr[i] >= p1 })
       ansL[qi] = cL
       // count segments with l > pts[cnt-1]
       pk := pts[cnt-1]
       idxR := sort.Search(n, func(i int) bool { return lArr[i] > pk })
       ansR[qi] = n - idxR
       // gaps
       for i := 0; i+1 < cnt; i++ {
           gaps = append(gaps, GapQuery{x: pts[i], y: pts[i+1], idx: qi})
       }
   }

   // process gap queries offline
   // sort segments by l descending
   sort.Slice(segs, func(i, j int) bool { return segs[i].l > segs[j].l })
   // sort gaps by x descending
   sort.Slice(gaps, func(i, j int) bool { return gaps[i].x > gaps[j].x })
   // BIT over r coordinates [1..maxR]
   const maxC = 1000005
   bit := NewBIT(maxC)
   si := 0
   for _, g := range gaps {
       // add segments with l > g.x
       for si < n && segs[si].l > g.x {
           r := segs[si].r
           if r >= 1 && r <= maxC {
               bit.Add(r, 1)
           }
           si++
       }
       // count segments with r < g.y => r <= g.y-1
       lim := g.y - 1
       if lim > 0 {
           cnt := bit.Sum(lim)
           gapSum[g.idx] += cnt
       }
   }

   // output answers
   for i := 0; i < m; i++ {
       res := n - ansL[i] - ansR[i] - gapSum[i]
       fmt.Fprintln(writer, res)
   }
}
