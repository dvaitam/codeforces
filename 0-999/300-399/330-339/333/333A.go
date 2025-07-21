package main

import (
   "fmt"
)

func main() {
   var n uint64
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   m := n
   // find least significant non-zero digit in base-3: count trailing zeros
   var p uint64
   for m%3 == 0 {
       m /= 3
       p++
   }
   // denom = 3^(p+1)
   denom := uint64(1)
   for i := uint64(0); i <= p; i++ {
       denom *= 3
   }
   // minimum coins needed = ceil(n / denom)
   ans := (n + denom - 1) / denom
   fmt.Println(ans)
}
