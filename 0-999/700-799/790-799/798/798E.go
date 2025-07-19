package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n+2)
   b := make([]int, n+2)
   ans := make([]int, n+2)
   for i := 1; i <= n; i++ {
       var ai int
       fmt.Fscan(reader, &ai)
       a[i] = ai
       if ai != -1 {
           b[ai] = i
       }
   }
   for i := 1; i <= n; i++ {
       if a[i] == -1 {
           a[i] = n + 1
       }
       if b[i] == 0 {
           b[i] = n + 1
       }
   }
   // segment tree arrays
   size := 4 * (n + 1)
   stVal := make([]int, size)
   stIdx := make([]int, size)
   // update position idx: set value b[idx]
   var update func(pos, l, r, idx int)
   update = func(pos, l, r, idx int) {
       if l == r {
           stVal[pos] = b[l]
           stIdx[pos] = l
           return
       }
       mid := (l + r) >> 1
       if idx <= mid {
           update(pos<<1, l, mid, idx)
       } else {
           update(pos<<1|1, mid+1, r, idx)
       }
       // merge children
       lc, rc := pos<<1, pos<<1|1
       if stVal[lc] > stVal[rc] || (stVal[lc] == stVal[rc] && stIdx[lc] > stIdx[rc]) {
           stVal[pos], stIdx[pos] = stVal[lc], stIdx[lc]
       } else {
           stVal[pos], stIdx[pos] = stVal[rc], stIdx[rc]
       }
   }
   // query max (value, idx) in [ql, qr]
   var query func(pos, l, r, ql, qr int) (int, int)
   query = func(pos, l, r, ql, qr int) (int, int) {
       if ql > r || qr < l {
           return 0, 0
       }
       if ql <= l && r <= qr {
           return stVal[pos], stIdx[pos]
       }
       mid := (l + r) >> 1
       lv, li := query(pos<<1, l, mid, ql, qr)
       rv, ri := query(pos<<1|1, mid+1, r, ql, qr)
       if lv > rv || (lv == rv && li > ri) {
           return lv, li
       }
       return rv, ri
   }
   // build tree
   for i := 1; i <= n; i++ {
       update(1, 1, n, i)
   }
   cnt := 0
   type frame struct{ u, stage, v int }
   stack := make([]frame, 0, n)
   // iterative DFS for topu
   topu := func(start int) {
       stack = append(stack, frame{u: start, stage: 0})
       for len(stack) > 0 {
           fr := &stack[len(stack)-1]
           u := fr.u
           switch fr.stage {
           case 0:
               v := b[u]
               fr.v = v
               b[u] = 0
               update(1, 1, n, u)
               if v != n+1 && b[v] != 0 {
                   fr.stage = 1
                   stack = append(stack, frame{u: v, stage: 0})
               } else {
                   fr.stage = 2
               }
           case 1:
               fr.stage = 2
           case 2:
               av := a[u] - 1
               var kv, ki int
               if av >= 1 {
                   kv, ki = query(1, 1, n, 1, av)
               }
               if kv > u {
                   fr.stage = 2
                   stack = append(stack, frame{u: ki, stage: 0})
               } else {
                   cnt++
                   ans[u] = cnt
                   stack = stack[:len(stack)-1]
               }
           }
       }
   }
   for i := n; i >= 1; i-- {
       if ans[i] == 0 {
           topu(i)
       }
   }
   // output answers
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, ans[i])
   }
}
