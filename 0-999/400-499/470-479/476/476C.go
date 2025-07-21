package main

import (
   "fmt"
)

const MOD int64 = 1000000007

func main() {
   var a, b int64
   if _, err := fmt.Scan(&a, &b); err != nil {
       return
   }
   inv2 := (MOD + 1) / 2
   aMod := a % MOD
   bMod := b % MOD

   // sum_k = 1 + 2 + ... + a = a*(a+1)/2 mod MOD
   sumK := aMod * ((aMod + 1) % MOD) % MOD * inv2 % MOD
   // t = b * sumK + a
   t := (bMod*sumK%MOD + aMod) % MOD
   // sum_r = sum_{r=1..b-1} r = (b-1)*b/2 mod MOD
   sumR := ((bMod - 1 + MOD) % MOD) * bMod % MOD * inv2 % MOD
   // answer = t * sumR mod MOD
   ans := t * sumR % MOD
   fmt.Println(ans)
}
