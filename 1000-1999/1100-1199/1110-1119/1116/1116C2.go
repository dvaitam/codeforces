package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// main reads multiple test cases and for each computes the maximum number of days
// the hero can be powered given initial powers and available power-ups (k).
func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for tc := 0; tc < t; tc++ {
       var n int
       var k int64
       fmt.Fscan(reader, &n, &k)
       p := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &p[i])
       }
       sort.Slice(p, func(i, j int) bool { return p[i] < p[j] })
       // binary search on days
       left, right := 0, n
       for left < right {
           mid := (left + right + 1) / 2
           var cost int64
           // consider mid days, use largest mid powers
           for i := 1; i <= mid; i++ {
               idx := n - mid + i - 1
               need := int64(i) - p[idx]
               if need > 0 {
                   cost += need
                   if cost > k {
                       break
                   }
               }
           }
           if cost <= k {
               left = mid
           } else {
               right = mid - 1
           }
       }
       fmt.Fprintln(writer, left)
   }
}
