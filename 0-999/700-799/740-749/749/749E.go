package main

import "fmt"

// Fenwick tree for int64 values
type Fenwick struct {
   n    int
   tree []int64
}

// NewFenwick creates a Fenwick tree of size n (1-indexed)
func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, tree: make([]int64, n+1)}
}

// Update adds v at index i
func (f *Fenwick) Update(i int, v int64) {
   for ; i <= f.n; i += i & -i {
       f.tree[i] += v
   }
}

// Query returns sum of [1..i]
func (f *Fenwick) Query(i int) int64 {
   var s int64
   for ; i > 0; i -= i & -i {
       s += f.tree[i]
   }
   return s
}

func main() {
   var N int
   if _, err := fmt.Scan(&N); err != nil {
       return
   }
   a := make([]int, N)
   for i := 0; i < N; i++ {
       fmt.Scan(&a[i])
   }

   fenCnt := NewFenwick(N)
   fenSum := NewFenwick(N)
   var good, totgood int64
   for i, ai := range a {
       if ai > 1 {
           good += fenCnt.Query(ai - 1)
           totgood += fenSum.Query(ai-1) * int64(N-i)
       }
       // update trees
       fenCnt.Update(ai, 1)
       fenSum.Update(ai, int64(i+1))
   }

   n := float64(N)
   dv := n * (n + 1.0)
   v := (n*(n-1.0)/2.0 - float64(good)) * dv
   f := (n - 1.0) * n * (n + 1.0) * (n + 2.0) / 24.0
   t := float64(totgood*2) - f
   res := (v + t) / dv
   fmt.Printf("%.16f\n", res)
}
