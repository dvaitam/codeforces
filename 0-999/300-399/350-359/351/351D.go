package main

import (
   "bufio"
   "fmt"
   "os"
)

// Fenwick tree for sum queries and point updates
type Fenwick struct {
   n    int
   tree []int
}

func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, tree: make([]int, n+1)}
}

// Add v at position i (1-based index)
func (f *Fenwick) Add(i, v int) {
   for x := i; x <= f.n; x += x & -x {
       f.tree[x] += v
   }
}

// Sum returns prefix sum up to i (1-based)
func (f *Fenwick) Sum(i int) int {
   s := 0
   for x := i; x > 0; x -= x & -x {
       s += f.tree[x]
   }
   return s
}

// Query returns sum in [l, r]
func (f *Fenwick) Query(l, r int) int {
   if l > r {
       return 0
   }
   return f.Sum(r) - f.Sum(l-1)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var m int
   if _, err := fmt.Fscan(in, &m); err != nil {
       return
   }
   b := make([]int, m+1)
   for i := 1; i <= m; i++ {
       fmt.Fscan(in, &b[i])
   }
   var q int
   fmt.Fscan(in, &q)
   // Queries grouped by right endpoint
   qs := make([][]struct{ l, idx int }, m+1)
   for i := 0; i < q; i++ {
       var l, r int
       fmt.Fscan(in, &l, &r)
       qs[r] = append(qs[r], struct{ l, idx int }{l, i})
   }
   // Prepare answer array
   ans := make([]int, q)
   // Fenwick tree over positions 1..m
   fw := NewFenwick(m)
   // Track last occurrence of each value
   const maxVal = 100000
   last := make([]int, maxVal+1)

   // Process prefixes
   for i := 1; i <= m; i++ {
       val := b[i]
       if last[val] != 0 {
           fw.Add(last[val], -1)
       }
       fw.Add(i, 1)
       last[val] = i
       // Answer queries ending at i
       for _, qu := range qs[i] {
           ans[qu.idx] = fw.Query(qu.l, i)
       }
   }
   // Output answers
   for i := 0; i < q; i++ {
       fmt.Fprintln(out, ans[i])
   }
}
