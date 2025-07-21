package main

import (
   "bufio"
   "fmt"
   "os"
)

// Node stores segment tree data for a segment
type Node struct {
   c4, c7, b47, b74, lazy int32
}

var (
   n, m int
   s    string
   tree []Node
)

func max(a, b int32) int32 {
   if a > b {
       return a
   }
   return b
}

// build constructs the tree over s[l..r] at position pos
func build(pos, l, r int) {
   if l == r {
       if s[l] == '4' {
           tree[pos].c4 = 1
       } else {
           tree[pos].c7 = 1
       }
       tree[pos].b47 = 1
       tree[pos].b74 = 1
       return
   }
   mid := (l + r) >> 1
   build(pos<<1, l, mid)
   build(pos<<1|1, mid+1, r)
   pull(pos)
}

// pull updates node pos from its children
func pull(pos int) {
   left := tree[pos<<1]
   right := tree[pos<<1|1]
   tree[pos].c4 = left.c4 + right.c4
   tree[pos].c7 = left.c7 + right.c7
   // longest non-decreasing subsequence: 4..7
   tree[pos].b47 = max(left.c4+right.b47, left.b47+right.c7)
   // longest non-increasing subsequence: 7..4
   tree[pos].b74 = max(left.c7+right.b74, left.b74+right.c4)
}

// applyFlip toggles 4<->7 in node pos
func applyFlip(pos int) {
   t := tree[pos]
   t.c4, t.c7 = t.c7, t.c4
   t.b47, t.b74 = t.b74, t.b47
   t.lazy ^= 1
   tree[pos] = t
}

// push propagates lazy flag to children
func push(pos int) {
   if tree[pos].lazy != 0 {
       applyFlip(pos << 1)
       applyFlip(pos<<1 | 1)
       tree[pos].lazy = 0
   }
}

// update flips digits in [ql..qr]
func update(pos, l, r, ql, qr int) {
   if ql > r || qr < l {
       return
   }
   if ql <= l && r <= qr {
       applyFlip(pos)
       return
   }
   push(pos)
   mid := (l + r) >> 1
   if ql <= mid {
       update(pos<<1, l, mid, ql, qr)
   }
   if qr > mid {
       update(pos<<1|1, mid+1, r, ql, qr)
   }
   pull(pos)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n, &m)
   fmt.Fscan(reader, &s)
   tree = make([]Node, 4*n+5)
   build(1, 0, n-1)
   for i := 0; i < m; i++ {
       var cmd string
       fmt.Fscan(reader, &cmd)
       if cmd[0] == 's' { // switch
           var l, r int
           fmt.Fscan(reader, &l, &r)
           update(1, 0, n-1, l-1, r-1)
       } else { // count
           fmt.Fprintln(writer, tree[1].b47)
       }
   }
}
