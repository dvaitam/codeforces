package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
)

// segment tree for range add and max query
type segTree struct {
   n    int
   mx   []int
   add  []int
}

func newSegTree(n int) *segTree {
   size := 1
   for size < n {
       size <<= 1
   }
   mx := make([]int, size<<1)
   add := make([]int, size<<1)
   return &segTree{n: size, mx: mx, add: add}
}

func (st *segTree) apply(p, v int) {
   st.mx[p] += v
   st.add[p] += v
}

func (st *segTree) push(p int) {
   if st.add[p] != 0 {
       st.apply(p<<1, st.add[p])
       st.apply(p<<1|1, st.add[p])
       st.add[p] = 0
   }
}

func (st *segTree) pull(p int) {
   if st.mx[p<<1] > st.mx[p<<1|1] {
       st.mx[p] = st.mx[p<<1]
   } else {
       st.mx[p] = st.mx[p<<1|1]
   }
}

// add v to [l, r]
func (st *segTree) updateRange(l, r, v int) {
   st.updateRangeRec(1, 0, st.n-1, l, r, v)
}

func (st *segTree) updateRangeRec(p, lo, hi, l, r, v int) {
   if l > hi || r < lo {
       return
   }
   if l <= lo && hi <= r {
       st.apply(p, v)
       return
   }
   st.push(p)
   mid := (lo + hi) >> 1
   st.updateRangeRec(p<<1, lo, mid, l, r, v)
   st.updateRangeRec(p<<1|1, mid+1, hi, l, r, v)
   st.pull(p)
}

type event struct {
   u, delta int
   l, r     int
}

func main() {
   br := bufio.NewReader(os.Stdin)
   var n, r int
   fmt.Fscan(br, &n, &r)
   pts := make([][2]int, n)
   for i := 0; i < n; i++ {
       var x, y int
       fmt.Fscan(br, &x, &y)
       pts[i][0], pts[i][1] = x, y
   }
   events := make([]event, 0, n*2)
   vs := make([]int, 0, n*2)
   for i := 0; i < n; i++ {
       x, y := pts[i][0], pts[i][1]
       u := x + y
       v := x - y
       l := v - r
       rr := v + r
       events = append(events, event{u - r, 1, l, rr})
       events = append(events, event{u + r, -1, l, rr})
       vs = append(vs, l, rr)
   }
   sort.Ints(vs)
   uni := vs[:1]
   for i := 1; i < len(vs); i++ {
       if vs[i] != vs[i-1] {
           uni = append(uni, vs[i])
       }
   }
   // map v to index
   idx := func(v int) int {
       i := sort.SearchInts(uni, v)
       return i
   }
   for i := range events {
       events[i].l = idx(events[i].l)
       events[i].r = idx(events[i].r)
   }
   sort.Slice(events, func(i, j int) bool {
       if events[i].u != events[j].u {
           return events[i].u < events[j].u
       }
       return events[i].delta > events[j].delta
   })
   st := newSegTree(len(uni))
   ans := 0
   for _, e := range events {
       st.updateRange(e.l, e.r, e.delta)
       if st.mx[1] > ans {
           ans = st.mx[1]
       }
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   w.WriteString(strconv.Itoa(ans))
}
