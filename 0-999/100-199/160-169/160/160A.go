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
   coins := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&coins[i])
   }
   sort.Sort(sort.Reverse(sort.IntSlice(coins)))
   total := 0
   for _, v := range coins {
       total += v
   }
   sum := 0
   for i, v := range coins {
       sum += v
       if sum > total - sum {
           fmt.Println(i + 1)
           return
       }
   }
}
