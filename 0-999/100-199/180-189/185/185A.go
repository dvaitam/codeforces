package main

import (
   "fmt"
)

const mod int64 = 1000000007

// fast exponentiation: compute (base^exp) % mod
func powMod(base, exp int64) int64 {
   res := int64(1)
   base %= mod
   for exp > 0 {
       if exp&1 == 1 {
           res = (res * base) % mod
       }
       base = (base * base) % mod
       exp >>= 1
   }
   return res
}

func main() {
   var n int64
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   // U(n) = (4^n + 2^n) / 2 mod mod
   p4 := powMod(4, n)
   p2 := powMod(2, n)
   // multiply by inverse of 2 modulo mod
   inv2 := (mod + 1) / 2
   ans := ( (p4 + p2) % mod * inv2 ) % mod
   fmt.Println(ans)
}
