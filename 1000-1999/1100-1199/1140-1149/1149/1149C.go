package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 1e18

// Node holds segment tree information
type Node struct {
   mn, mx    int   // minimum and maximum prefix sum
   l1, r1    int   // helper values for cross intervals
   best      int   // max diameter within segment
   lazy      int   // pending addition
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, q int
   fmt.Fscan(reader, &n, &q)
   var s string
   fmt.Fscan(reader, &s)
   bs := []byte(s)
   m := len(bs)
   // prefix sums s[0..m]
   ps := make([]int, m+1)
   for i := 1; i <= m; i++ {
       if bs[i-1] == '(' {
           ps[i] = ps[i-1] + 1
       } else {
           ps[i] = ps[i-1] - 1
       }
   }
   // build segment tree of size N = m+1
   N := m + 1
   size := 1
   for size < N {
       size <<= 1
   }
   tree := make([]Node, size<<1)
   negInf := int(-INF)
   // initialize leaves
   for i := 0; i < size; i++ {
       idx := i + size
       if i < N {
           tree[idx].mn = ps[i]
           tree[idx].mx = ps[i]
           tree[idx].l1 = negInf
           tree[idx].r1 = negInf
           tree[idx].best = 0
       } else {
           // empty beyond N
           tree[idx].mn = int(1e9)
           tree[idx].mx = int(-1e9)
           tree[idx].l1 = negInf
           tree[idx].r1 = negInf
           tree[idx].best = 0
       }
   }
   // build internal nodes
   for i := size - 1; i >= 1; i-- {
       tree[i] = merge(tree[i<<1], tree[i<<1|1])
   }
   // output initial diameter
   fmt.Fprintln(writer, tree[1].best)
   // process queries
   for k := 0; k < q; k++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       l, r := a, b
       if l > r {
           l, r = r, l
       }
       // delta on [l..r-1]
       var delta int
       if bs[l-1] == '(' {
           delta = -2
       } else {
           delta = 2
       }
       update(1, 0, size-1, l, r-1, delta, tree)
       // swap brackets
       bs[l-1], bs[r-1] = bs[r-1], bs[l-1]
       // print updated diameter
       fmt.Fprintln(writer, tree[1].best)
   }
}

// merge two nodes
func merge(a, b Node) Node {
   c := Node{}
   // min and max
   c.mn = a.mn
   if b.mn < c.mn {
       c.mn = b.mn
   }
   c.mx = a.mx
   if b.mx > c.mx {
       c.mx = b.mx
   }
   // l1: max(a.l1, a.mx-2*b.mn, b.l1)
   c.l1 = a.l1
   if v := a.mx - 2*b.mn; v > c.l1 {
       c.l1 = v
   }
   if b.l1 > c.l1 {
       c.l1 = b.l1
   }
   // r1: max(b.r1, b.mx-2*a.mn, a.r1)
   c.r1 = b.r1
   if v := b.mx - 2*a.mn; v > c.r1 {
       c.r1 = v
   }
   if a.r1 > c.r1 {
       c.r1 = a.r1
   }
   // best: max(a.best, b.best, a.l1+b.mx, a.mx+b.r1)
   c.best = a.best
   if b.best > c.best {
       c.best = b.best
   }
   if v := a.l1 + b.mx; v > c.best {
       c.best = v
   }
   if v := a.mx + b.r1; v > c.best {
       c.best = v
   }
   return c
}

// apply lazy tag
func applyTag(n *Node, d int) {
   n.mn += d
   n.mx += d
   n.l1 -= d
   n.r1 -= d
   n.lazy += d
}

// push propagates tag to children
func push(p int, tree []Node) {
   if tree[p].lazy != 0 {
       d := tree[p].lazy
       applyTag(&tree[p<<1], d)
       applyTag(&tree[p<<1|1], d)
       tree[p].lazy = 0
   }
}

// update range [ql..qr] by d
func update(p, l, r, ql, qr, d int, tree []Node) {
   if ql > r || qr < l {
       return
   }
   if ql <= l && r <= qr {
       applyTag(&tree[p], d)
       return
   }
   push(p, tree)
   mid := (l + r) >> 1
   update(p<<1, l, mid, ql, qr, d, tree)
   update(p<<1|1, mid+1, r, ql, qr, d, tree)
   tree[p] = merge(tree[p<<1], tree[p<<1|1])
}
