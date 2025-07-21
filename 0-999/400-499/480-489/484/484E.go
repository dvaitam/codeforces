package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// BIT for range prefix max query and point update
type BIT struct {
   n int
   t []int
}

func NewBIT(n int) *BIT {
   return &BIT{n: n, t: make([]int, n+1)}
}

// Update position i (1-indexed) with value v: t[i] = max(t[i], v)
func (b *BIT) Update(i, v int) {
   for ; i <= b.n; i += i & -i {
       if b.t[i] < v {
           b.t[i] = v
       }
   }
}

// Query max on prefix [1..i]
func (b *BIT) Query(i int) int {
   res := 0
   for ; i > 0; i -= i & -i {
       if b.t[i] > res {
           res = b.t[i]
       }
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   h := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &h[i])
   }
   // compute L[i]
   L := make([]int, n+1)
   stack := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       for len(stack) > 0 && h[stack[len(stack)-1]] >= h[i] {
           stack = stack[:len(stack)-1]
       }
       if len(stack) == 0 {
           L[i] = 0
       } else {
           L[i] = stack[len(stack)-1]
       }
       stack = append(stack, i)
   }
   // compute R[i]
   R := make([]int, n+1)
   stack = stack[:0]
   for i := n; i >= 1; i-- {
       for len(stack) > 0 && h[stack[len(stack)-1]] >= h[i] {
           stack = stack[:len(stack)-1]
       }
       if len(stack) == 0 {
           R[i] = n + 1
       } else {
           R[i] = stack[len(stack)-1]
       }
       stack = append(stack, i)
   }
   // bars: a = L[i], b = R[i], h = h[i]
   type Bar struct{ a, b, h int }
   bars := make([]Bar, n)
   for i := 1; i <= n; i++ {
       bars[i-1] = Bar{L[i], R[i], h[i]}
   }
   sort.Slice(bars, func(i, j int) bool { return bars[i].a < bars[j].a })

   var m int
   fmt.Fscan(in, &m)
   type Query struct{ t1, t2, idx int }
   qs := make([]Query, m)
   res := make([]int, m)
   for i := 0; i < m; i++ {
       var l, r, w int
       fmt.Fscan(in, &l, &r, &w)
       t1 := r - w
       t2 := l + w
       qs[i] = Query{t1, t2, i}
   }
   sort.Slice(qs, func(i, j int) bool { return qs[i].t1 < qs[j].t1 })
   // BIT on transformed b: idx = (n+2 - b)
   size := n + 2
   bit := NewBIT(size)
   bi := 0
   for _, q := range qs {
       // insert bars with a <= t1
       for bi < n && bars[bi].a <= q.t1 {
           bval := bars[bi].b
           idx := size - bval
           // BIT is 1-indexed, idx from 0..size-1, shift by +1
           bit.Update(idx+1, bars[bi].h)
           bi++
       }
       // query for b >= t2 => idx <= size - t2
       k := size - q.t2
       if k < 0 {
           res[q.idx] = 0
       } else {
           res[q.idx] = bit.Query(k + 1)
       }
   }
   for i := 0; i < m; i++ {
       fmt.Fprintln(out, res[i])
   }
}
