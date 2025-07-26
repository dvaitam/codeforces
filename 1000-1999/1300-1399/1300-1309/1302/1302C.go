package main

import (
   "bufio"
   "fmt"
   "os"
)

// Fenwick tree for range sum, 1-based index
type Fenwick struct {
   n    int
   tree []int64
}

// NewFenwick creates a Fenwick tree of size n
func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, tree: make([]int64, n+1)}
}

// Update adds delta at position i
func (f *Fenwick) Update(i int, delta int64) {
   for ; i <= f.n; i += i & -i {
       f.tree[i] += delta
   }
}

// PrefixSum returns sum of [1..i]
func (f *Fenwick) PrefixSum(i int) int64 {
   var s int64
   for ; i > 0; i -= i & -i {
       s += f.tree[i]
   }
   return s
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for T > 0 {
       T--
       var n, q int
       fmt.Fscan(reader, &n, &q)
       fw := NewFenwick(n)
       arr := make([]int64, n+1)
       for i := 0; i < q; i++ {
           var typ int
           fmt.Fscan(reader, &typ)
           if typ == 1 {
               var x int
               var y int64
               fmt.Fscan(reader, &x, &y)
               diff := y - arr[x]
               arr[x] = y
               fw.Update(x, diff)
           } else {
               var l, r int
               fmt.Fscan(reader, &l, &r)
               sum := fw.PrefixSum(r) - fw.PrefixSum(l-1)
               fmt.Fprintln(writer, sum)
           }
       }
   }
}
