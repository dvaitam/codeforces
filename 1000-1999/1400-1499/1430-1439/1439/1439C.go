package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 1 << 60

// Segment tree supporting range chmax (assign a[i] = max(a[i], x)), range min queries, and range sum queries.
// Implemented via segment tree beats for chmax (lower bound updates).
type Node struct {
   sum    int64
   min1   int
   min2   int
   cntMin int
}

type SegTree struct {
   n    int
   tree []Node
}

func NewSegTree(a []int) *SegTree {
   n := len(a) - 1 // a is 1-indexed with len = n+1
   st := &SegTree{n: n, tree: make([]Node, 4*(n+1))}
   st.build(1, 1, n, a)
   return st
}

func (st *SegTree) build(u, l, r int, a []int) {
   if l == r {
       v := a[l]
       st.tree[u] = Node{sum: int64(v), min1: v, min2: INF, cntMin: 1}
       return
   }
   m := (l + r) >> 1
   st.build(u<<1, l, m, a)
   st.build(u<<1|1, m+1, r, a)
   st.pushUp(u)
}

func (st *SegTree) pushUp(u int) {
   L, R := st.tree[u<<1], st.tree[u<<1|1]
   // sum
   st.tree[u].sum = L.sum + R.sum
   // min1, min2, cntMin
   if L.min1 < R.min1 {
       st.tree[u].min1 = L.min1
       st.tree[u].cntMin = L.cntMin
       st.tree[u].min2 = min(L.min2, R.min1)
   } else if L.min1 > R.min1 {
       st.tree[u].min1 = R.min1
       st.tree[u].cntMin = R.cntMin
       st.tree[u].min2 = min(L.min1, R.min2)
   } else {
       st.tree[u].min1 = L.min1
       st.tree[u].cntMin = L.cntMin + R.cntMin
       st.tree[u].min2 = min(L.min2, R.min2)
   }
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

// apply chmax update to node u with value x (raise min1 to x)
func (st *SegTree) applyChmax(u, x int) {
   node := &st.tree[u]
   if node.min1 >= x {
       return
   }
   // only elements equal to min1 are increased to x
   node.sum += int64(x-node.min1) * int64(node.cntMin)
   node.min1 = x
}

func (st *SegTree) pushDown(u int) {
   for i := 0; i < 2; i++ {
       v := u<<1 | i
       if st.tree[v].min1 < st.tree[u].min1 {
           st.applyChmax(v, st.tree[u].min1)
       }
   }
}

// rangeChmax on [ql, qr]: a[i] = max(a[i], x)
func (st *SegTree) rangeChmax(u, l, r, ql, qr, x int) {
   if r < ql || qr < l || st.tree[u].min1 >= x {
       return
   }
   if ql <= l && r <= qr && st.tree[u].min2 > x {
       st.applyChmax(u, x)
       return
   }
   st.pushDown(u)
   m := (l + r) >> 1
   st.rangeChmax(u<<1, l, m, ql, qr, x)
   st.rangeChmax(u<<1|1, m+1, r, ql, qr, x)
   st.pushUp(u)
}

// findFirstLE finds the first index i >= ql such that a[i] <= y, or returns n+1
func (st *SegTree) findFirstLE(u, l, r, ql int, y int64) int {
   if r < ql || int64(st.tree[u].min1) > y {
       return st.n + 1
   }
   if l == r {
       return l
   }
   st.pushDown(u)
   m := (l + r) >> 1
   if m >= ql {
       res := st.findFirstLE(u<<1, l, m, ql, y)
       if res <= st.n {
           return res
       }
       return st.findFirstLE(u<<1|1, m+1, r, ql, y)
   }
   return st.findFirstLE(u<<1|1, m+1, r, ql, y)
}

// queryPrefix returns the last index r >= ql such that sum(a[ql..r]) <= *y; *y is decreased by sum(a[ql..r])
func (st *SegTree) queryPrefix(u, l, r, ql int, y *int64) int {
   if r < ql || *y < int64(st.tree[u].min1) {
       return ql - 1
   }
   if ql <= l && st.tree[u].sum <= *y {
       *y -= st.tree[u].sum
       return r
   }
   if l == r {
       *y -= st.tree[u].sum
       return l
   }
   st.pushDown(u)
   m := (l + r) >> 1
   var res int
   if m >= ql {
       res = st.queryPrefix(u<<1, l, m, ql, y)
       if res < m {
           return res
       }
       return st.queryPrefix(u<<1|1, m+1, r, ql, y)
   }
   return st.queryPrefix(u<<1|1, m+1, r, ql, y)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, q int
   fmt.Fscan(reader, &n, &q)
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   st := NewSegTree(a)
   for i := 0; i < q; i++ {
       var t, x int
       var y int64
       fmt.Fscan(reader, &t, &x, &y)
       if t == 1 {
           st.rangeChmax(1, 1, n, 1, x, int(y))
       } else {
           cnt := 0
           pos := x
           for pos <= n {
               p := st.findFirstLE(1, 1, n, pos, y)
               if p > n {
                   break
               }
               r := st.queryPrefix(1, 1, n, p, &y)
               cnt += r - p + 1
               pos = r + 1
           }
           fmt.Fprintln(writer, cnt)
       }
   }
}
