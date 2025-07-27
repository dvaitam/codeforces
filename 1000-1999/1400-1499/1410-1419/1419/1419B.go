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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   // Precompute nice staircase costs S = n*(n+1)/2 for n = 2^k -1
   const limit = int64(1e18)
   var S []int64
   for k := 1; ; k++ {
       n := (1 << uint(k)) - 1
       // compute cost = n*(n+1)/2, stop if exceeds limit
       cost := int64(n) * int64(n+1) / 2
       if cost > limit {
           break
       }
       S = append(S, cost)
   }
   // prefix sums
   P := make([]int64, len(S))
   for i, v := range S {
       if i == 0 {
           P[i] = v
       } else {
           P[i] = P[i-1] + v
       }
   }

   for ti := 0; ti < t; ti++ {
       var x int64
       fmt.Fscan(reader, &x)
       // find max k such that P[k] <= x
       idx := sort.Search(len(P), func(i int) bool { return P[i] > x })
       // idx is count of values <= x
       fmt.Fprintln(writer, idx)
   }
}
