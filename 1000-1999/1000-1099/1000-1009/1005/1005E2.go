package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Fenwick tree for sum queries
type Fenwick struct {
   n    int
   tree []int
}

func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, tree: make([]int, n+1)}
}

// Add value v at index i (1-based)
func (f *Fenwick) Add(i, v int) {
   for x := i; x <= f.n; x += x & -x {
       f.tree[x] += v
   }
}

// Sum returns sum of [1..i]
func (f *Fenwick) Sum(i int) int {
   s := 0
   for x := i; x > 0; x -= x & -x {
       s += f.tree[x]
   }
   return s
}

// countGE returns number of subarrays where median >= k
func countGE(a []int, k int) int64 {
   n := len(a)
   // prefix sums
   ps := make([]int, n+1)
   for i := 1; i <= n; i++ {
       if a[i-1] >= k {
           ps[i] = ps[i-1] + 1
       } else {
           ps[i] = ps[i-1] - 1
       }
   }
   // compress values
   vals := make([]int, n+1)
   copy(vals, ps)
   sort.Ints(vals)
   uniq := vals[:0]
   for _, v := range vals {
       if len(uniq) == 0 || uniq[len(uniq)-1] != v {
           uniq = append(uniq, v)
       }
   }
   // Fenwick on uniq size
   fw := NewFenwick(len(uniq))
   var res int64
   // include ps[0]
   // iterate i from 0..n
   for i := 0; i <= n; i++ {
       // find index of ps[i]
       idx := sort.SearchInts(uniq, ps[i]) + 1 // 1-based
       // count number of previous ps[j] < ps[i]
       if idx > 1 {
           res += int64(fw.Sum(idx - 1))
       }
       // add current ps[i]
       fw.Add(idx, 1)
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // result = count subarrays with median >= m minus >= m+1
   ans := countGE(a, m) - countGE(a, m+1)
   fmt.Fprintln(writer, ans)
}
