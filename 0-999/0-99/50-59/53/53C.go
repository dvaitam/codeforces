package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   res := make([]int, 0, n)
   // Construct permutation: 1, n, 2, n-1, ...
   for i := 1; i <= n/2; i++ {
       res = append(res, i)
       res = append(res, n-i+1)
   }
   if n%2 == 1 {
       res = append(res, (n+1)/2)
   }
   // Output
   for i, v := range res {
       if i > 0 {
           fmt.Print(" ")
       }
       fmt.Print(v)
   }
}
