package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const INF = 1000000000

// Segment tree node arrays will be 1-indexed by node
type SegTree struct {
   n     int
   arr   []int
   lucky []int
   mn    []int
   sum   []int
   lazy  []int
}

func NewSegTree(n int, arr []int, lucky []int) *SegTree {
   size := 4 * n
   st := &SegTree{
       n:     n,
       arr:   make([]int, n),
       lucky: lucky,
       mn:    make([]int, size),
       sum:   make([]int, size),
       lazy:  make([]int, size),
   }
   copy(st.arr, arr)
   st.build(1, 0, n-1)
   return st
}

// build initializes mn and sum
func (st *SegTree) build(node, l, r int) {
   if l == r {
       d := st.nextDist(st.arr[l])
       st.mn[node] = d
       if d == 0 {
           st.sum[node] = 1
       }
       return
   }
   mid := (l + r) >> 1
   st.build(node<<1, l, mid)
   st.build(node<<1|1, mid+1, r)
   st.pull(node)
}

func (st *SegTree) pull(node int) {
   l, r := node<<1, node<<1|1
   if st.mn[l] < st.mn[r] {
       st.mn[node] = st.mn[l]
   } else {
       st.mn[node] = st.mn[r]
   }
   st.sum[node] = st.sum[l] + st.sum[r]
}

func (st *SegTree) push(node int) {
   if st.lazy[node] != 0 {
       d := st.lazy[node]
       for _, c := range []int{node << 1, node<<1 | 1} {
           st.lazy[c] += d
           st.mn[c] -= d
       }
       st.lazy[node] = 0
   }
}

// nextDist returns distance to next lucky >= x, or INF
func (st *SegTree) nextDist(x int) int {
   idx := sort.Search(len(st.lucky), func(i int) bool { return st.lucky[i] >= x })
   if idx < len(st.lucky) {
       return st.lucky[idx] - x
   }
   return INF
}

// Update adds d to arr in [ql,qr]
func (st *SegTree) Update(node, l, r, ql, qr, d int) {
   if r < ql || l > qr {
       return
   }
   if ql <= l && r <= qr && st.mn[node] > d {
       st.mn[node] -= d
       st.lazy[node] += d
       return
   }
   if l == r {
       // apply pending
       st.arr[l] += st.lazy[node]
       st.lazy[node] = 0
       // apply this update
       st.arr[l] += d
       // recompute
       dd := st.nextDist(st.arr[l])
       st.mn[node] = dd
       if dd == 0 {
           st.sum[node] = 1
       } else {
           st.sum[node] = 0
       }
       return
   }
   st.push(node)
   mid := (l + r) >> 1
   st.Update(node<<1, l, mid, ql, qr, d)
   st.Update(node<<1|1, mid+1, r, ql, qr, d)
   st.pull(node)
}

// Query returns count of lucky numbers in [ql,qr]
func (st *SegTree) Query(node, l, r, ql, qr int) int {
   if r < ql || l > qr {
       return 0
   }
   if ql <= l && r <= qr {
       return st.sum[node]
   }
   st.push(node)
   mid := (l + r) >> 1
   return st.Query(node<<1, l, mid, ql, qr) + st.Query(node<<1|1, mid+1, r, ql, qr)
}

// generate all lucky numbers <= max
func genLuck(cur, maxv int, out *[]int) {
   if cur > maxv {
       return
   }
   if cur > 0 {
       *out = append(*out, cur)
   }
   genLuck(cur*10+4, maxv, out)
   genLuck(cur*10+7, maxv, out)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   arr := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &arr[i])
   }
   var lucky []int
   genLuck(0, 10000, &lucky)
   sort.Ints(lucky)
   st := NewSegTree(n, arr, lucky)
   for i := 0; i < m; i++ {
       var op string
       fmt.Fscan(reader, &op)
       if op == "add" {
           var l, r, d int
           fmt.Fscan(reader, &l, &r, &d)
           st.Update(1, 0, n-1, l-1, r-1, d)
       } else {
           var l, r int
           fmt.Fscan(reader, &l, &r)
           res := st.Query(1, 0, n-1, l-1, r-1)
           fmt.Fprintln(writer, res)
       }
   }
}
