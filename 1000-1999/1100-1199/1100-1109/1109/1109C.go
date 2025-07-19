package main

import (
   "bufio"
   "fmt"
   "os"
)

const L = 1000000009

// Range represents a segment's aggregated data
type Range struct {
   l, r, f, c int
   low, su   int64
   has       bool
}

// merge combines two Ranges
func merge(a, b Range) Range {
   var ret Range
   ret.l = a.l
   ret.r = b.r
   if a.has {
       ret.f = a.f
   } else {
       ret.f = b.f
   }
   if b.has {
       ret.c = b.c
   } else {
       ret.c = a.c
   }
   mid := int64(b.f - b.l) * int64(a.c)
   lowb := b.low
   if lowb > 0 {
       lowb = 0
   }
   if a.low < a.su+mid+lowb {
       ret.low = a.low
   } else {
       ret.low = a.su + mid + lowb
   }
   ret.su = a.su + mid + b.su
   ret.has = a.has || b.has
   return ret
}

// newRangeSeg creates a default segment range [a,b]
func newRangeSeg(a, b int) Range {
   return Range{l: a, r: b, f: b + 1}
}

// newRangePoint creates a point update range [a,b] with value v
func newRangePoint(a, b, v int) Range {
   vv := int64(v)
   return Range{l: a, r: b, f: a, c: v, low: vv, su: vv, has: true}
}

// newRangeVal creates initial query Range from v
func newRangeVal(v int) Range {
   vv := int64(v)
   return Range{f: 1, low: vv, su: vv}
}

// Node of the dynamic segment tree
type Node struct {
   left, right *Node
   v           Range
}

var (
   root *Node
   qq   Range
   ok   bool
   ans  float64
)

func newNode(l, r int) *Node {
   return &Node{v: newRangeSeg(l, r)}
}

// upd performs point update at p with value v
func upd(x **Node, l, r, p, v int) {
   if *x == nil {
       *x = newNode(l, r)
   }
   cur := *x
   if l < r {
       mid := (l + r) >> 1
       if p <= mid {
           upd(&cur.left, l, mid, p, v)
       } else {
           upd(&cur.right, mid+1, r, p, v)
       }
       var leftR, rightR Range
       if cur.left != nil {
           leftR = cur.left.v
       } else {
           leftR = newRangeSeg(l, mid)
       }
       if cur.right != nil {
           rightR = cur.right.v
       } else {
           rightR = newRangeSeg(mid+1, r)
       }
       cur.v = merge(leftR, rightR)
   } else {
       cur.v = newRangePoint(p, p, v)
   }
}

// clr clears the point at p
func clr(x **Node, l, r, p int) {
   if *x == nil {
       *x = newNode(l, r)
   }
   cur := *x
   if l < r {
       mid := (l + r) >> 1
       if p <= mid {
           clr(&cur.left, l, mid, p)
       } else {
           clr(&cur.right, mid+1, r, p)
       }
       var leftR, rightR Range
       if cur.left != nil {
           leftR = cur.left.v
       } else {
           leftR = newRangeSeg(l, mid)
       }
       if cur.right != nil {
           rightR = cur.right.v
       } else {
           rightR = newRangeSeg(mid+1, r)
       }
       cur.v = merge(leftR, rightR)
   } else {
       cur.v = newRangeSeg(p, p)
   }
}

// solve checks the first position in [l,r] where condition meets
func solve(x *Node, l, r int) {
   if ok {
       return
   }
   if x == nil || !x.v.has {
       su := qq.su
       vv := int64(qq.c)
       tt := int64(r - l + 1)
       if su+vv*tt <= 0 {
           ok = true
           ans = float64(l) + float64(su)/float64(-vv)
       } else {
           qq = merge(qq, newRangeSeg(l, r))
       }
       return
   }
   f := x.v.f
   ins := int64(qq.c) * int64(minInt(r+1, f)-l)
   lowb := x.v.low
   if lowb > 0 {
       lowb = 0
   }
   if qq.su+ins+lowb > 0 {
       qq = merge(qq, x.v)
       return
   }
   if qq.su+ins <= 0 {
       ok = true
       ans = float64(l) + float64(qq.su)/float64(-qq.c)
       return
   }
   if f > r {
       return
   }
   if l == r {
       if qq.su+int64(x.v.c) <= 0 {
           ok = true
           ans = float64(l) + float64(qq.su)/float64(-x.v.c)
       }
   } else {
       mid := (l + r) >> 1
       solve(x.left, l, mid)
       if ok {
           return
       }
       solve(x.right, mid+1, r)
   }
}

// qry traverses the tree to query [ql,qr]
func qry(x *Node, l, r, ql, qr int) {
   if ok {
       return
   }
   if x == nil || (l == ql && r == qr) {
       solve(x, ql, qr)
       return
   }
   mid := (l + r) >> 1
   if ql <= mid {
       qry(x.left, l, mid, ql, minInt(mid, qr))
   }
   if qr > mid {
       qry(x.right, mid+1, r, maxInt(mid+1, ql), qr)
   }
}

func minInt(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func maxInt(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var Q int
   fmt.Fscan(reader, &Q)
   for i := 0; i < Q; i++ {
       var op int
       fmt.Fscan(reader, &op)
       if op == 1 {
           var t, s int
           fmt.Fscan(reader, &t, &s)
           upd(&root, 0, L, t, s)
       } else if op == 2 {
           var t int
           fmt.Fscan(reader, &t)
           clr(&root, 0, L, t)
       } else {
           var l, r, v int
           fmt.Fscan(reader, &l, &r, &v)
           if v == 0 {
               fmt.Fprintln(writer, l)
               continue
           }
           qq = newRangeVal(v)
           ok = false
           qry(root, 0, L, l, r-1)
           if ok && ans >= float64(l) && ans <= float64(r) {
               fmt.Fprintf(writer, "%.6f\n", ans)
           } else {
               fmt.Fprintln(writer, -1)
           }
       }
   }
}
