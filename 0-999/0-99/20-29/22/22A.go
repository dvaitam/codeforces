package main

import (
   "fmt"
   "sort"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   nums := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&nums[i])
   }
   // Collect distinct values
   uniq := make(map[int]struct{})
   for _, v := range nums {
       uniq[v] = struct{}{}
   }
   // Need at least two distinct values
   if len(uniq) < 2 {
       fmt.Println("NO")
       return
   }
   // Sort unique values
   vals := make([]int, 0, len(uniq))
   for v := range uniq {
       vals = append(vals, v)
   }
   sort.Ints(vals)
   // Second order statistic is the second smallest
   fmt.Println(vals[1])
}
