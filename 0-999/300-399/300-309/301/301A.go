package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   m := 2*n - 1
   sumAbs := 0
   negCount := 0
   minAbs := int(1e9)
   for i := 0; i < m; i++ {
       var x int
       fmt.Scan(&x)
       if x < 0 {
           negCount++
       }
       ax := abs(x)
       sumAbs += ax
       if ax < minAbs {
           minAbs = ax
       }
   }
   // If n is odd, any configuration reachable; else ensure even number of flips
   if n%2 == 1 || negCount%2 == 0 {
       fmt.Println(sumAbs)
   } else {
       fmt.Println(sumAbs - 2*minAbs)
   }
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}
