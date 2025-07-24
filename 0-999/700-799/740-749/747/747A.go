package main

import (
   "fmt"
   "math"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   // start from the integer part of sqrt(n) and go down to find factors
   a := int(math.Sqrt(float64(n)))
   for a > 0 {
       if n%a == 0 {
           b := n / a
           fmt.Println(a, b)
           return
       }
       a--
   }
}
