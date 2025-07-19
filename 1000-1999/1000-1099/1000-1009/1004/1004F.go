package main

import (
   "bufio"
   "fmt"
   "os"
)

// pair holds OR value and position
type pair struct { fst, scd int }

// segment stores prefix/suffix OR streaks and count of subarrays with OR >= X
type segment struct {
   vL  []pair
   vR  []pair
   sum int64
}

var (
   n, m, X int
   arr      []int
   tree     []segment
)

// merge combines two segments into one
func merge(a, b segment) segment {
   var c segment
   // build c.vL
   c.vL = make([]pair, len(a.vL))
   copy(c.vL, a.vL)
   for _, p := range b.vL {
       pre := c.vL[len(c.vL)-1]
       newVal := pre.fst | p.fst
       if newVal > pre.fst {
           c.vL = append(c.vL, pair{newVal, p.scd})
       }
   }
   // build c.vR
   c.vR = make([]pair, len(b.vR))
   copy(c.vR, b.vR)
   for _, p := range a.vR {
       pre := c.vR[len(c.vR)-1]
       newVal := pre.fst | p.fst
       if newVal > pre.fst {
           c.vR = append(c.vR, pair{newVal, p.scd})
       }
   }
   // sum parts
   c.sum = a.sum + b.sum
   // count cross subarrays
   j := 0
   for i := len(a.vR) - 1; i >= 0; i-- {
       // advance j until OR >= X
       for j < len(b.vL) && ((a.vR[i].fst | b.vL[j].fst) < X) {
           j++
       }
       if j < len(b.vL) {
           // left count
           leftEnd := a.vR[i].scd
           var leftStart int
           if i+1 < len(a.vR) {
               leftStart = a.vR[i+1].scd
           } else {
               leftStart = a.vL[0].scd - 1
           }
           cntL := leftEnd - leftStart
           // right count
           cntR := b.vR[0].scd - b.vL[j].scd + 1
           c.sum += int64(cntL) * int64(cntR)
       }
   }
   return c
}

// build tree
func build(rt, l, r int) {
   if l == r {
       val := arr[l]
       tree[rt].vL = []pair{{val, l}}
       tree[rt].vR = []pair{{val, l}}
       if val >= X {
           tree[rt].sum = 1
       }
   } else {
       m := (l + r) >> 1
       lc, rc := rt<<1, rt<<1|1
       build(lc, l, m)
       build(rc, m+1, r)
       tree[rt] = merge(tree[lc], tree[rc])
   }
}

// update point
func update(rt, l, r, pos, val int) {
   if l == r {
       tree[rt].vL = []pair{{val, l}}
       tree[rt].vR = []pair{{val, l}}
       if val >= X {
           tree[rt].sum = 1
       } else {
           tree[rt].sum = 0
       }
   } else {
       m := (l + r) >> 1
       lc, rc := rt<<1, rt<<1|1
       if pos <= m {
           update(lc, l, m, pos, val)
       } else {
           update(rc, m+1, r, pos, val)
       }
       tree[rt] = merge(tree[lc], tree[rc])
   }
}

// query returns segment covering [ql,qr]
func query(rt, l, r, ql, qr int) segment {
   if ql <= l && r <= qr {
       return tree[rt]
   }
   m := (l + r) >> 1
   lc, rc := rt<<1, rt<<1|1
   if qr <= m {
       return query(lc, l, m, ql, qr)
   }
   if ql > m {
       return query(rc, m+1, r, ql, qr)
   }
   left := query(lc, l, m, ql, qr)
   right := query(rc, m+1, r, ql, qr)
   return merge(left, right)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n, &m, &X)
   arr = make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &arr[i])
   }
   tree = make([]segment, 4*(n+1))
   build(1, 1, n)
   for k := 0; k < m; k++ {
       var typ, x, y int
       fmt.Fscan(in, &typ, &x, &y)
       if typ == 1 {
           update(1, 1, n, x, y)
       } else {
           res := query(1, 1, n, x, y)
           fmt.Fprintln(out, res.sum)
       }
   }
}
