package main

import (
   "fmt"
   "math"
)

func main() {
   var n int64
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   var m int64 = 1
   for {
       x := n*m + 1
       s := int64(math.Sqrt(float64(x)))
       if s*s == x {
           fmt.Println(m)
           break
       }
       m++
   }
}
