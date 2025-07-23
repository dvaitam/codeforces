package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// segment tree for range add and range min
type SegTree struct {
   n    int
   minv []int
   lazy []int
}

func NewSegTree(n int) *SegTree {
   size := 1
   for size < n {
       size <<= 1
   }
   minv := make([]int, size*2)
   lazy := make([]int, size*2)
   return &SegTree{n: n, minv: minv, lazy: lazy}
}

func (st *SegTree) apply(p, v int) {
   st.minv[p] += v
   st.lazy[p] += v
}

func (st *SegTree) push(p int) {
   if st.lazy[p] != 0 {
       st.apply(p*2, st.lazy[p])
       st.apply(p*2+1, st.lazy[p])
       st.lazy[p] = 0
   }
}

func (st *SegTree) pull(p int) {
   if st.minv[p*2] < st.minv[p*2+1] {
       st.minv[p] = st.minv[p*2]
   } else {
       st.minv[p] = st.minv[p*2+1]
   }
}

func (st *SegTree) updateRange(p, l, r, ql, qr, v int) {
   if ql > r || qr < l {
       return
   }
   if ql <= l && r <= qr {
       st.apply(p, v)
       return
   }
   st.push(p)
   m := (l + r) >> 1
   st.updateRange(p*2, l, m, ql, qr, v)
   st.updateRange(p*2+1, m+1, r, ql, qr, v)
   st.pull(p)
}

func (st *SegTree) queryMin(p, l, r, ql, qr int) int {
   if ql > r || qr < l {
       // return large
       return 1<<60
   }
   if ql <= l && r <= qr {
       return st.minv[p]
   }
   st.push(p)
   m := (l + r) >> 1
   a := st.queryMin(p*2, l, m, ql, qr)
   b := st.queryMin(p*2+1, m+1, r, ql, qr)
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   var s int
   fmt.Fscan(reader, &s)
   type Query struct{ l, r int; x int }
   qs := make([]Query, s)
   for i := 0; i < s; i++ {
       fmt.Fscan(reader, &qs[i].l, &qs[i].r, &qs[i].x)
       qs[i].l--
       qs[i].r--
   }

   // prepare events: for each a_j, window range [l_j,r_j]
   W := n - m + 1
   if W < 1 {
       W = 1
   }
   type P struct{ v, idx int }
   ps := make([]P, n)
   for j := 0; j < n; j++ {
       ps[j] = P{v: a[j], idx: j}
   }
   sort.Slice(ps, func(i, j int) bool { return ps[i].v < ps[j].v })

   st := NewSegTree(W)
   // initial all zeros
   p := 0
   curQ := 0
   prevAns := 0
   for i := 0; i < s; i++ {
       qi := qs[i].x ^ prevAns
       // adjust events
       if qi > curQ {
           for p < n && ps[p].v < qi {
               j := ps[p].idx
               l := j - m + 1
               if l < 0 {
                   l = 0
               }
               r := j
               if r >= W {
                   r = W - 1
               }
               if l <= r {
                   st.updateRange(1, 0, st.n-1, l, r, 1)
               }
               p++
           }
       } else if qi < curQ {
           for p > 0 && ps[p-1].v >= qi {
               p--
               j := ps[p].idx
               l := j - m + 1
               if l < 0 {
                   l = 0
               }
               r := j
               if r >= W {
                   r = W - 1
               }
               if l <= r {
                   st.updateRange(1, 0, st.n-1, l, r, -1)
               }
           }
       }
       curQ = qi
       // query range [l_i, r_i] on windows
       // windows v corresponds to start positions [0..W-1], query.l..query.r
       // ensure r_i <= W-1
       li := qs[i].l
       ri := qs[i].r
       if ri >= W {
           ri = W - 1
       }
       ans := st.queryMin(1, 0, st.n-1, li, ri)
       if ans < 0 {
           ans = 0
       }
       fmt.Fprint(writer, ans)
       if i+1 < s {
           fmt.Fprint(writer, " ")
       }
       prevAns = ans
   }
   writer.WriteByte('\n')
}
