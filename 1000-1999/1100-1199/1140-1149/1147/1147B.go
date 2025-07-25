package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   adj := make([][]int, n)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       a--
       b--
       adj[a] = append(adj[a], b)
       adj[b] = append(adj[b], a)
   }
   // compute hash for neighbor distance lists
   H := make([]uint64, n)
   const P = 1000003
   for i := 0; i < n; i++ {
       neigh := adj[i]
       dists := make([]int, len(neigh))
       for j, v := range neigh {
           d := v - i
           if d < 0 {
               d += n
           }
           dists[j] = d
       }
       sort.Ints(dists)
       var h uint64 = 146527 // arbitrary seed
       for _, d := range dists {
           h = h*P + uint64(d)
       }
       // incorporate length to distinguish
       h = h*P + uint64(len(dists))
       H[i] = h
   }
   // find divisors of n (excluding n)
   divisors := make([]int, 0, 16)
   for d := 1; d*d <= n; d++ {
       if n%d == 0 {
           divisors = append(divisors, d)
           if d != n/d {
               divisors = append(divisors, n/d)
           }
       }
   }
   // check each proper divisor
   for _, d := range divisors {
       if d == n {
           continue
       }
       ok := true
       for i := 0; i < n; i++ {
           j := i + d
           if j >= n {
               j -= n
           }
           if H[i] != H[j] {
               ok = false
               break
           }
       }
       if ok {
           fmt.Fprintln(out, "Yes")
           return
       }
   }
   fmt.Fprintln(out, "No")
}
