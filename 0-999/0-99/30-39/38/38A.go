package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   d := make([]int, n+1)
   for i := 1; i < n; i++ {
       fmt.Scan(&d[i])
   }
   var a, b int
   fmt.Scan(&a, &b)
   sum := 0
   for i := a; i < b; i++ {
       sum += d[i]
   }
   fmt.Println(sum)
}
