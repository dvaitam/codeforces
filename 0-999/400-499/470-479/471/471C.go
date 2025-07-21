package main

import (
   "fmt"
)

func main() {
   var n int64
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   var count int64
   // h floors, minimal cards needed: h*(3*h+1)/2
   for h := int64(1); ; h++ {
       // compute minimal cards for h floors
       req := h*(3*h+1) / 2
       if req > n {
           break
       }
       // total cards n must satisfy (n+h) divisible by 3
       if (n+h)%3 == 0 {
           count++
       }
   }
   fmt.Println(count)
}
