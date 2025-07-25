package main

import (
   "bufio"
   "fmt"
   "os"
)

// Segment tree for range max and first position with value>=x
type SegTree struct {
   n    int
   size int
   mx   []int64
}

func NewSegTree(a []int64) *SegTree {
   n := len(a)
   size := 1
   for size < n {
       size <<= 1
   }
   mx := make([]int64, 2*size)
   for i := 0; i < n; i++ {
       mx[size+i] = a[i]
   }
   for i := size - 1; i > 0; i-- {
       if mx[2*i] > mx[2*i+1] {
           mx[i] = mx[2*i]
       } else {
           mx[i] = mx[2*i+1]
       }
   }
   return &SegTree{n, size, mx}
}
// query first index >= l with value >= x, within [l,r]
func (st *SegTree) queryFirst(l, r int, x int64) int {
   return st.qf(1, 0, st.size-1, l, r, x)
}
func (st *SegTree) qf(node, nl, nr, ql, qr int, x int64) int {
   if nr < ql || nl > qr || st.mx[node] < x {
       return -1
   }
   if nl == nr {
       if nl < st.n {
           return nl
       }
       return -1
   }
   mid := (nl + nr) >> 1
   // search left
   res := st.qf(2*node, nl, mid, ql, qr, x)
   if res != -1 {
       return res
   }
   return st.qf(2*node+1, mid+1, nr, ql, qr, x)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m int
   var p int64
   fmt.Fscan(in, &n, &m, &p)
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   pref := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       pref[i] = pref[i-1] + a[i]
   }
   // build segment tree on pref
   st := NewSegTree(pref)
   for qi := 0; qi < m; qi++ {
       var l, r int
       fmt.Fscan(in, &l, &r)
       base := pref[l-1]
       thr := base + p
       pos := l
       cnt := int64(0)
       for pos <= r {
           idx := st.queryFirst(pos, r, thr)
           if idx == -1 {
               break
           }
           cnt++
           pos = idx + 1
           thr += p
       }
       raw := pref[r] - base
       res := raw - cnt*p
       fmt.Fprintln(out, res)
   }
}
