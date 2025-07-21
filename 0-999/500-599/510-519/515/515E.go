package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = int64(4e18)

// Node stores segment tree node information
type Node struct {
   mxL, mxR, best int64
}

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

// merge two nodes
func merge(a, b Node) Node {
   res := Node{}
   res.mxL = max(a.mxL, b.mxL)
   res.mxR = max(a.mxR, b.mxR)
   res.best = max(max(a.best, b.best), a.mxL + b.mxR)
   return res
}

var seg []Node
var Larr, Rarr []int64

// build segment tree on [l, r] at index p
func build(p, l, r int) {
   if l == r {
       seg[p] = Node{mxL: Larr[l], mxR: Rarr[l], best: -INF}
       return
   }
   m := (l + r) >> 1
   lc, rc := p<<1, p<<1|1
   build(lc, l, m)
   build(rc, m+1, r)
   seg[p] = merge(seg[lc], seg[rc])
}

// query segment [ql, qr] on node p covering [l, r]
func query(p, l, r, ql, qr int) Node {
   if qr < l || r < ql {
       return Node{mxL: -INF, mxR: -INF, best: -INF}
   }
   if ql <= l && r <= qr {
       return seg[p]
   }
   m := (l + r) >> 1
   left := query(p<<1, l, m, ql, qr)
   right := query(p<<1|1, m+1, r, ql, qr)
   return merge(left, right)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   d := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &d[i])
   }
   h := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &h[i])
   }
   // build prefix distances on doubled cycle
   N := 2 * n
   C := make([]int64, N+2)
   for i := 1; i <= N; i++ {
       di := d[(i-1)%n+1]
       C[i] = C[i-1] + di
   }
   Larr = make([]int64, N+2)
   Rarr = make([]int64, N+2)
   for i := 1; i <= N; i++ {
       hi := h[(i-1)%n+1]
       Larr[i] = 2*hi - C[i-1]
       Rarr[i] = 2*hi + C[i-1]
   }
   // build segment tree
   seg = make([]Node, 4*(N+1))
   build(1, 1, N)

   // process queries
   for qi := 0; qi < m; qi++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       var l, r int
       if a <= b {
           l = b + 1
           r = a - 1 + n
       } else {
           l = b + 1
           r = a - 1
       }
       // query on [l, r]
       res := query(1, 1, N, l, r)
       fmt.Fprintln(writer, res.best)
   }
}
