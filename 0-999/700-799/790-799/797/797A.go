package main

import (
   "fmt"
)

func main() {
   var n, k int
   if _, err := fmt.Scan(&n, &k); err != nil {
       return
   }
   temp := k
   if k == 1 {
       fmt.Println(n)
       return
   }
   var factors []int
   i := 2
   for i <= (n+1)/2+1 {
       if n%i == 0 {
           k--
           factors = append(factors, i)
           n /= i
           if n == 1 {
               n = 0
           }
           i = 2
       } else {
           i++
           continue
       }
       if n == 0 {
           break
       }
       if k == 1 {
           k--
           factors = append(factors, n)
           break
       }
   }
   if len(factors) != temp || k != 0 {
       fmt.Println(-1)
       return
   }
   for idx, v := range factors {
       if idx > 0 {
           fmt.Print(" ")
       }
       fmt.Print(v)
   }
   fmt.Println()
}
