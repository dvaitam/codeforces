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
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // If only one segment, cost is full range
   if k <= 1 {
       fmt.Println(a[n-1] - a[0])
       return
   }
   // Compute consecutive differences
   diffs := make([]int64, 0, n-1)
   for i := 1; i < n; i++ {
       diffs = append(diffs, a[i]-a[i-1])
   }
   // Sort descending
   sort.Slice(diffs, func(i, j int) bool {
       return diffs[i] > diffs[j]
   })
   // Select k-1 largest gaps to cut
   var sumMaxGaps int64
   cuts := k - 1
   for i := 0; i < cuts && i < len(diffs); i++ {
       sumMaxGaps += diffs[i]
   }
   // Total cost = full range minus saved gaps
   cost := (a[n-1] - a[0]) - sumMaxGaps
   fmt.Println(cost)
}
