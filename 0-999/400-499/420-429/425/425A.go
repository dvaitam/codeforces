package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // prefix sums
   prefix := make([]int, n+1)
   for i := 0; i < n; i++ {
       prefix[i+1] = prefix[i] + a[i]
   }
   const INF = 1e18
   // use int64 for sums
   var maxSum int64 = -INF
   // iterate over all subarrays [l, r]
   for l := 0; l < n; l++ {
       for r := l; r < n; r++ {
           // base sum
           base := int64(prefix[r+1] - prefix[l])
           // collect inside and outside elements
           inside := make([]int, 0, r-l+1)
           outside := make([]int, 0, n-(r-l+1))
           for i := 0; i < n; i++ {
               if i >= l && i <= r {
                   inside = append(inside, a[i])
               } else {
                   outside = append(outside, a[i])
               }
           }
           sort.Ints(inside)             // ascending
           sort.Sort(sort.Reverse(sort.IntSlice(outside))) // descending
           // compute best with up to k swaps
           curSum := base
           if curSum > maxSum {
               maxSum = curSum
           }
           // number of possible swaps
           maxSwap := k
           if len(inside) < maxSwap {
               maxSwap = len(inside)
           }
           if len(outside) < maxSwap {
               maxSwap = len(outside)
           }
           // try swapping
           for t := 0; t < maxSwap; t++ {
               if outside[t] > inside[t] {
                   curSum += int64(outside[t] - inside[t])
                   if curSum > maxSum {
                       maxSum = curSum
                   }
               } else {
                   break
               }
           }
       }
   }
   // output result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, maxSum)
}
