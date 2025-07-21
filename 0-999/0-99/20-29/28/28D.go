package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

const INF = math.MaxInt64 / 4

type segTree struct {
   n    int
   min  []int64
   idx  []int
   lazy []int64
}

func newSegTree(a []int64) *segTree {
   n := len(a) - 1
   size := 4 * (n + 1)
   st := &segTree{n: n, min: make([]int64, size), idx: make([]int, size), lazy: make([]int64, size)}
   var build func(node, l, r int)
   build = func(node, l, r int) {
       if l == r {
           st.min[node] = a[l]
           st.idx[node] = l
       } else {
           m := (l + r) >> 1
           build(node<<1, l, m)
           build(node<<1|1, m+1, r)
           if st.min[node<<1] <= st.min[node<<1|1] {
               st.min[node] = st.min[node<<1]
               st.idx[node] = st.idx[node<<1]
           } else {
               st.min[node] = st.min[node<<1|1]
               st.idx[node] = st.idx[node<<1|1]
           }
       }
   }
   build(1, 1, n)
   return st
}

func (st *segTree) apply(node int, v int64) {
   st.min[node] += v
   st.lazy[node] += v
}

func (st *segTree) push(node int) {
   if st.lazy[node] != 0 {
       st.apply(node<<1, st.lazy[node])
       st.apply(node<<1|1, st.lazy[node])
       st.lazy[node] = 0
   }
}

func (st *segTree) pull(node int) {
   if st.min[node<<1] <= st.min[node<<1|1] {
       st.min[node] = st.min[node<<1]
       st.idx[node] = st.idx[node<<1]
   } else {
       st.min[node] = st.min[node<<1|1]
       st.idx[node] = st.idx[node<<1|1]
   }
}

// range add v on [L,R]
func (st *segTree) updateRange(node, l, r, L, R int, v int64) {
   if L > r || R < l {
       return
   }
   if L <= l && r <= R {
       st.apply(node, v)
       return
   }
   st.push(node)
   m := (l + r) >> 1
   st.updateRange(node<<1, l, m, L, R, v)
   st.updateRange(node<<1|1, m+1, r, L, R, v)
   st.pull(node)
}

// set position p to INF
func (st *segTree) updatePoint(node, l, r, p int) {
   if l == r {
       st.min[node] = INF
       st.lazy[node] = 0
       return
   }
   st.push(node)
   m := (l + r) >> 1
   if p <= m {
       st.updatePoint(node<<1, l, m, p)
   } else {
       st.updatePoint(node<<1|1, m+1, r, p)
   }
   st.pull(node)
}

// get global min and its index
func (st *segTree) queryMin() (int64, int) {
   return st.min[1], st.idx[1]
}

func main() {
   rdr := bufio.NewReader(os.Stdin)
   wrt := bufio.NewWriter(os.Stdout)
   defer wrt.Flush()
   var n int
   fmt.Fscan(rdr, &n)
   v := make([]int, n+1)
   c := make([]int64, n+1)
   l := make([]int64, n+1)
   r := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(rdr, &v[i], &c[i], &l[i], &r[i])
   }
   // prefix sums
   pref := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       pref[i] = pref[i-1] + c[i]
   }
   // suffix sums
   suff := make([]int64, n+2)
   for i := n; i >= 1; i-- {
       suff[i] = suff[i+1] + c[i]
   }
   // compute initial slacks
   slack1 := make([]int64, n+1)
   slack2 := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       slack1[i] = pref[i-1] - l[i]
       slack2[i] = suff[i+1] - r[i]
   }
   st1 := newSegTree(slack1)
   st2 := newSegTree(slack2)
   alive := make([]bool, n+1)
   for i := 1; i <= n; i++ {
       alive[i] = true
   }
   // iterative removal
   for {
       mn1, i1 := st1.queryMin()
       mn2, i2 := st2.queryMin()
       if mn1 >= 0 && mn2 >= 0 {
           break
       }
       var k int
       if mn1 < mn2 {
           k = i1
       } else {
           k = i2
       }
       if !alive[k] {
           // already removed, guard
           if mn1 < 0 {
               st1.updatePoint(1, 1, n, i1)
           }
           if mn2 < 0 {
               st2.updatePoint(1, 1, n, i2)
           }
           continue
       }
       // remove k
       alive[k] = false
       // update range effects
       if k+1 <= n {
           st1.updateRange(1, 1, n, k+1, n, -c[k])
       }
       if k-1 >= 1 {
           st2.updateRange(1, 1, n, 1, k-1, -c[k])
       }
       // mark k as removed
       st1.updatePoint(1, 1, n, k)
       st2.updatePoint(1, 1, n, k)
   }
   // collect and output
   var res []int
   for i := 1; i <= n; i++ {
       if alive[i] {
           res = append(res, i)
       }
   }
   fmt.Fprintln(wrt, len(res))
   if len(res) > 0 {
       for i, id := range res {
           if i > 0 {
               fmt.Fprint(wrt, ' ')
           }
           fmt.Fprint(wrt, id)
       }
       fmt.Fprintln(wrt)
   }
}
