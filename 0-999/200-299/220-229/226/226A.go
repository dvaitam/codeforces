package main

import (
   "fmt"
)

// modPow computes (base^exp) % mod using binary exponentiation.
func modPow(base, exp, mod int64) int64 {
   result := int64(1)
   base %= mod
   for exp > 0 {
       if exp&1 == 1 {
           result = (result * base) % mod
       }
       base = (base * base) % mod
       exp >>= 1
   }
   return result
}

func main() {
   var n, m int64
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   // answer is (3^n - 1) mod m
   if m == 1 {
       fmt.Println(0)
       return
   }
   p := modPow(3, n, m)
   ans := p - 1
   if ans < 0 {
       ans += m
   }
   fmt.Println(ans)
}
