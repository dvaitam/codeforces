package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Fenwick tree (BIT) for sum of int64
type Fenwick struct {
   n    int
   tree []int64
}

func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, tree: make([]int64, n+1)}
}

// Add v at position i (1-based)
func (f *Fenwick) Add(i int, v int64) {
   for ; i <= f.n; i += i & -i {
       f.tree[i] += v
   }
}

// Sum returns sum of [1..i]
func (f *Fenwick) Sum(i int) int64 {
   var s int64
   for ; i > 0; i -= i & -i {
       s += f.tree[i]
   }
   return s
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   var t int64
   fmt.Fscan(in, &n, &t)
   arr := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &arr[i])
   }
   // prefix sums
   pref := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       pref[i] = pref[i-1] + arr[i-1]
   }
   // coordinate compression
   comp := make([]int64, n+1)
   copy(comp, pref)
   sort.Slice(comp, func(i, j int) bool { return comp[i] < comp[j] })
   // unique
   m := 0
   for i := 0; i <= n; i++ {
       if m == 0 || comp[i] != comp[m-1] {
           comp[m] = comp[i]
           m++
       }
   }
   // Fenwick tree size m
   fw := NewFenwick(m)
   var ans int64
   // initial prefix 0
   // find index of 0
   idx0 := sort.Search(m, func(i int) bool { return comp[i] >= 0 }) + 1
   fw.Add(idx0, 1)
   // iterate
   for i := 1; i <= n; i++ {
       // count of pref[j] <= pref[i] - t
       target := pref[i] - t
       ub := sort.Search(m, func(j int) bool { return comp[j] > target }) // first > target
       cntLE := fw.Sum(ub)
       // number of previous prefixes is i
       ans += int64(i) - cntLE
       // add current prefix
       idx := sort.Search(m, func(j int) bool { return comp[j] >= pref[i] }) + 1
       fw.Add(idx, 1)
   }
   fmt.Fprintln(out, ans)
}
