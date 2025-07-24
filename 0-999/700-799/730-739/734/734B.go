package main

import (
   "fmt"
)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   var k2, k3, k5, k6 int64
   _, _ = fmt.Scan(&k2, &k3, &k5, &k6)
   // number of 256 integers we can form
   x := min(k2, min(k5, k6))
   k2 -= x
   // number of 32 integers we can form with remaining 2s and 3s
   y := min(k2, k3)
   result := x*256 + y*32
   fmt.Println(result)
}
