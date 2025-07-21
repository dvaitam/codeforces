package main

import (
   "fmt"
)

func main() {
   var n, b int
   if _, err := fmt.Scan(&n, &b); err != nil {
       return
   }
   a := make([]int, n)
   sum := 0
   for i := 0; i < n; i++ {
       fmt.Scan(&a[i])
       sum += a[i]
   }
   V := float64(sum+b) / float64(n)
   // Check feasibility: cannot remove, so final volume must be >= current
   for i := 0; i < n; i++ {
       if float64(a[i]) > V {
           fmt.Println(-1)
           return
       }
   }
   // Output additions
   for i := 0; i < n; i++ {
       c := V - float64(a[i])
       fmt.Printf("%.6f\n", c)
   }
}
