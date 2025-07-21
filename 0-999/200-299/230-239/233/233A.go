package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   if n%2 != 0 {
       fmt.Println(-1)
       return
   }
   p := make([]int, n)
   for i := 0; i < n; i += 2 {
       p[i] = i + 2
       p[i+1] = i + 1
   }
   for i := 0; i < n; i++ {
       if i > 0 {
           fmt.Print(" ")
       }
       fmt.Print(p[i])
   }
   fmt.Println()
}
