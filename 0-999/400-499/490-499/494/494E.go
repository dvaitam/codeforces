package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type event struct{
   x    int64
   y1, y2 int64
   delta int
}

// segment tree over y intervals storing covered length parity
type SegTree struct{
   l, r int
   cover int
   parity bool
   fullParity bool
   left, right *SegTree
}

func build(y []int64, l, r int) *SegTree {
   node := &SegTree{l: l, r: r}
   // compute full segment parity
   if l+1 == r {
       // leaf covers [y[l], y[r])
       node.fullParity = ((y[r] - y[l]) & 1) == 1
   } else {
       m := (l + r) >> 1
       node.left = build(y, l, m)
       node.right = build(y, m, r)
       node.fullParity = node.left.fullParity != node.right.fullParity
   }
   return node
}

func (st *SegTree) update(y []int64, ql, qr int, delta int) {
   if qr <= st.l || st.r <= ql {
       return
   }
   if ql <= st.l && st.r <= qr {
       st.cover += delta
   } else {
       st.left.update(y, ql, qr, delta)
       st.right.update(y, ql, qr, delta)
   }
   if st.cover > 0 {
       st.parity = st.fullParity
   } else if st.l+1 == st.r {
       st.parity = false
   } else {
       st.parity = st.left.parity != st.right.parity
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int64
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   events := make([]event, 0, m*2)
   ys := make([]int64, 0, m*2+2)
   for i := int64(0); i < m; i++ {
       var a, b, c, d int64
       fmt.Fscan(in, &a, &b, &c, &d)
       // add rectangle [a,c] x [b,d]
       events = append(events, event{x: a, y1: b, y2: d + 1, delta: 1})
       events = append(events, event{x: c + 1, y1: b, y2: d + 1, delta: -1})
       ys = append(ys, b, d+1)
   }
   if len(events) == 0 {
       fmt.Println("Malek")
       return
   }
   // compress y coordinates
   ys = append(ys, 1, n+1)
   sort.Slice(ys, func(i, j int) bool { return ys[i] < ys[j] })
   // unique
   uniq := make([]int64, 0, len(ys))
   for _, v := range ys {
       if len(uniq) == 0 || uniq[len(uniq)-1] != v {
           uniq = append(uniq, v)
       }
   }
   ys = uniq
   // map values
   yi := func(v int64) int {
       return sort.Search(len(ys), func(i int) bool { return ys[i] >= v })
   }
   // sort events by x
   sort.Slice(events, func(i, j int) bool { return events[i].x < events[j].x })
   // build segment tree
   st := build(ys, 0, len(ys)-1)
   var areaParity bool
   prevX := events[0].x
   i := 0
   for i < len(events) {
       x := events[i].x
       dx := x - prevX
       if dx & 1 == 1 {
           areaParity = areaParity != st.parity
       }
       // process all events at x
       for i < len(events) && events[i].x == x {
           l := yi(events[i].y1)
           r := yi(events[i].y2)
           st.update(ys, l, r, events[i].delta)
           i++
       }
       prevX = x
   }
   if areaParity {
       fmt.Println("Hamed")
   } else {
       fmt.Println("Malek")
   }
}
