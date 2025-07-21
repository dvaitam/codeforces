package main

import (
   "fmt"
)

// factor returns the exponents of primes 2,3,5 in x and the remaining residual
func factor(x int64) (e2, e3, e5 int, rem int64) {
   rem = x
   for rem%2 == 0 {
       rem /= 2
       e2++
   }
   for rem%3 == 0 {
       rem /= 3
       e3++
   }
   for rem%5 == 0 {
       rem /= 5
       e5++
   }
   return
}

func main() {
   var a, b int64
   if _, err := fmt.Scan(&a, &b); err != nil {
       return
   }
   a2, a3, a5, ra := factor(a)
   b2, b3, b5, rb := factor(b)
   if ra != rb {
       fmt.Println(-1)
       return
   }
   ops := abs(a2-b2) + abs(a3-b3) + abs(a5-b5)
   fmt.Println(ops)
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}
