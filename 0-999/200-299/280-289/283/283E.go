package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   fmt.Fscan(reader, &n, &k)
   skills := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &skills[i])
   }
   sorted := make([]int, n)
   copy(sorted, skills)
   sort.Ints(sorted)

   // difference arrays for flips coverage count and sum of l+r
   dt := make([]int64, n+3)
   dsum := make([]int64, n+3)

   for i := 0; i < k; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       // find interval [l0, r0] in sorted where skill in [a,b]
       l0 := sort.Search(n, func(i int) bool { return sorted[i] >= a })
       r0 := sort.Search(n, func(i int) bool { return sorted[i] > b }) - 1
       if l0 <= r0 {
           l := int64(l0 + 1)
           r := int64(r0 + 1)
           dt[l]++
           dt[r+1]--
           dsum[l] += l + r
           dsum[r+1] -= l + r
       }
   }

   // prefix sums to get t[u] and sum_lr[u]
   t := make([]int64, n+2)
   sumlr := make([]int64, n+2)
   var ct, cs int64
   for u := 1; u <= n; u++ {
       ct += dt[u]
       cs += dsum[u]
       t[u] = ct
       sumlr[u] = cs
   }

   // compute sum of transitive triples = sum_u C(outdeg[u],2)
   var trans int64
   for u := 1; u <= n; u++ {
       // initial outdeg = u-1; flips change by sumlr - 2*u*t
       out := int64(u-1) + sumlr[u] - 2*int64(u)*t[u]
       if out > 1 {
           trans += out * (out - 1) / 2
       }
   }
   // total triples
   nn := int64(n)
   total := nn * (nn - 1) * (nn - 2) / 6
   // cycles = total - transitive
   fmt.Fprint(writer, total-trans)
}
