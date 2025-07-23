package main

import (
   "bufio"
   "fmt"
   "os"
)

// segment tree supporting range add and querying leftmost/rightmost index where value equals y
type SegTree struct {
   n    int
   minv []int64
   maxv []int64
   lazy []int64
}

// NewSegTree builds a segment tree from the initial array a (1-based, length n+1)
func NewSegTree(a []int64) *SegTree {
   n := len(a) - 1
   size := 4 * n
   st := &SegTree{
       n:    n,
       minv: make([]int64, size),
       maxv: make([]int64, size),
       lazy: make([]int64, size),
   }
   var build func(node, l, r int)
   build = func(node, l, r int) {
       if l == r {
           st.minv[node] = a[l]
           st.maxv[node] = a[l]
           return
       }
       mid := (l + r) >> 1
       left, right := node<<1, node<<1|1
       build(left, l, mid)
       build(right, mid+1, r)
       // pull up
       if st.minv[left] < st.minv[right] {
           st.minv[node] = st.minv[left]
       } else {
           st.minv[node] = st.minv[right]
       }
       if st.maxv[left] > st.maxv[right] {
           st.maxv[node] = st.maxv[left]
       } else {
           st.maxv[node] = st.maxv[right]
       }
   }
   build(1, 1, n)
   return st
}

func (st *SegTree) apply(node int, v int64) {
   st.minv[node] += v
   st.maxv[node] += v
   st.lazy[node] += v
}

func (st *SegTree) push(node int) {
   if st.lazy[node] != 0 {
       v := st.lazy[node]
       st.apply(node<<1, v)
       st.apply(node<<1|1, v)
       st.lazy[node] = 0
   }
}

// Update adds x to range [ql, qr]
func (st *SegTree) Update(node, l, r, ql, qr int, x int64) {
   if ql <= l && r <= qr {
       st.apply(node, x)
       return
   }
   st.push(node)
   mid := (l + r) >> 1
   if ql <= mid {
       st.Update(node<<1, l, mid, ql, qr, x)
   }
   if qr > mid {
       st.Update(node<<1|1, mid+1, r, ql, qr, x)
   }
   left, right := node<<1, node<<1|1
   // pull up
   if st.minv[left] < st.minv[right] {
       st.minv[node] = st.minv[left]
   } else {
       st.minv[node] = st.minv[right]
   }
   if st.maxv[left] > st.maxv[right] {
       st.maxv[node] = st.maxv[left]
   } else {
       st.maxv[node] = st.maxv[right]
   }
}

const INF = 1<<60

// find leftmost index i in [l,r] where current value == y, return INF if none
func (st *SegTree) QueryLeft(node, l, r int, y int64) int {
   if st.minv[node] > y || st.maxv[node] < y {
       return int(INF)
   }
   if l == r {
       return l
   }
   st.push(node)
   mid := (l + r) >> 1
   res := st.QueryLeft(node<<1, l, mid, y)
   if res != int(INF) {
       return res
   }
   return st.QueryLeft(node<<1|1, mid+1, r, y)
}

// find rightmost index i in [l,r] where current value == y, return -INF if none
func (st *SegTree) QueryRight(node, l, r int, y int64) int {
   if st.minv[node] > y || st.maxv[node] < y {
       return -int(INF)
   }
   if l == r {
       return l
   }
   st.push(node)
   mid := (l + r) >> 1
   // try right first
   res := st.QueryRight(node<<1|1, mid+1, r, y)
   if res != -int(INF) {
       return res
   }
   return st.QueryRight(node<<1, l, mid, y)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, q int
   fmt.Fscan(reader, &n, &q)
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   st := NewSegTree(a)
   for ; q > 0; q-- {
       var t int
       fmt.Fscan(reader, &t)
       if t == 1 {
           var l, r int
           var x int64
           fmt.Fscan(reader, &l, &r, &x)
           st.Update(1, 1, n, l, r, x)
       } else if t == 2 {
           var y int64
           fmt.Fscan(reader, &y)
           if st.minv[1] > y || st.maxv[1] < y {
               fmt.Fprintln(writer, -1)
               continue
           }
           left := st.QueryLeft(1, 1, n, y)
           right := st.QueryRight(1, 1, n, y)
           if left > right {
               fmt.Fprintln(writer, -1)
           } else {
               fmt.Fprintln(writer, right-left)
           }
       }
   }
}
