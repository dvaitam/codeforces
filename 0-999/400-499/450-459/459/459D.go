package main

import (
   "bufio"
   "fmt"
   "os"
)

// Fenwick implements a Fenwick tree for sum queries.
type Fenwick struct {
   n   int
   bit []int
}

// NewFenwick returns a Fenwick tree of size n (1-indexed).
func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, bit: make([]int, n+1)}
}

// Add adds value v at index i.
func (f *Fenwick) Add(i, v int) {
   for ; i <= f.n; i += i & -i {
       f.bit[i] += v
   }
}

// Sum returns the prefix sum of [1..i].
func (f *Fenwick) Sum(i int) int {
   s := 0
   if i > f.n {
       i = f.n
   }
   for ; i > 0; i -= i & -i {
       s += f.bit[i]
   }
   return s
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   pre := make([]int, n)
   cnt := make(map[int]int)
   for i, x := range a {
       cnt[x]++
       pre[i] = cnt[x]
   }

   suf := make([]int, n)
   cnt = make(map[int]int)
   for i := n - 1; i >= 0; i-- {
       x := a[i]
       cnt[x]++
       suf[i] = cnt[x]
   }

   // Build Fenwick tree on suffix frequencies
   ft := NewFenwick(n)
   for i := 0; i < n; i++ {
       ft.Add(suf[i], 1)
   }

   var ans int64
   for i := 0; i < n; i++ {
       // remove current position from future queries
       ft.Add(suf[i], -1)
       v := pre[i]
       if v > 1 {
           ans += int64(ft.Sum(v - 1))
       }
   }
   fmt.Fprint(writer, ans)
}
