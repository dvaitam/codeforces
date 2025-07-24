package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Fenwick implements XOR Fenwick tree (1-based).
type Fenwick struct {
   n int
   t []int
}

// NewFenwick creates a Fenwick tree of size n.
func NewFenwick(n int) *Fenwick {
   return &Fenwick{n, make([]int, n+1)}
}

// Update applies XOR with v at position i.
func (f *Fenwick) Update(i, v int) {
   for ; i <= f.n; i += i & -i {
       f.t[i] ^= v
   }
}

// Query returns XOR of range [1..i].
func (f *Fenwick) Query(i int) int {
   res := 0
   for ; i > 0; i -= i & -i {
       res ^= f.t[i]
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // prefixXor for odd-frequency XOR
   prefixXor := make([]int, n+1)
   for i := 1; i <= n; i++ {
       prefixXor[i] = prefixXor[i-1] ^ a[i]
   }
   var m int
   fmt.Fscan(in, &m)
   // queries grouped by right endpoint
   type query struct{ l, idx int }
   qs := make([][]query, n+1)
   for i := 0; i < m; i++ {
       var l, r int
       fmt.Fscan(in, &l, &r)
       if r >= 1 && r <= n {
           qs[r] = append(qs[r], query{l, i})
       }
   }
   // compress values for last occurrence tracking
   vals := make([]int, n)
   for i := 1; i <= n; i++ {
       vals[i-1] = a[i]
   }
   sort.Ints(vals)
   u := 0
   for i := 0; i < n; i++ {
       if i == 0 || vals[i] != vals[i-1] {
           vals[u] = vals[i]
           u++
       }
   }
   vals = vals[:u]
   comp := make(map[int]int, u)
   for i, v := range vals {
       comp[v] = i
   }
   aComp := make([]int, n+1)
   for i := 1; i <= n; i++ {
       aComp[i] = comp[a[i]]
   }
   lastOcc := make([]int, u)
   bit := NewFenwick(n)
   ans := make([]int, m)
   // process
   for i := 1; i <= n; i++ {
       v := a[i]
       id := aComp[i]
       if p := lastOcc[id]; p != 0 {
           bit.Update(p, v)
       }
       bit.Update(i, v)
       lastOcc[id] = i
       for _, q := range qs[i] {
           totalDistinctXor := bit.Query(i) ^ bit.Query(q.l - 1)
           oddXor := prefixXor[i] ^ prefixXor[q.l-1]
           ans[q.idx] = totalDistinctXor ^ oddXor
       }
   }
   // output answers
   for i := 0; i < m; i++ {
       fmt.Fprintln(out, ans[i])
   }
}
