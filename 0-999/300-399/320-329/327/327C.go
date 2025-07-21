package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

// fast exponentiation: computes x^e mod MOD
func modPow(x, e int64) int64 {
   res := int64(1)
   x %= MOD
   for e > 0 {
       if e&1 == 1 {
           res = (res * x) % MOD
       }
       x = (x * x) % MOD
       e >>= 1
   }
   return res
}

// modular inverse of x mod MOD, MOD must be prime
func modInv(x int64) int64 {
   // Fermat's little theorem
   return modPow(x, MOD-2)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a string
   var k int64
   if _, err := fmt.Fscan(reader, &a); err != nil {
       return
   }
   fmt.Fscan(reader, &k)
   m := int64(len(a))
   // sum_p = sum of 2^p for positions p where a[p] == '0' or '5'
   var sumP int64 = 0
   var pow2 int64 = 1
   for i := int64(0); i < m; i++ {
       c := a[i]
       if c == '0' || c == '5' {
           sumP = (sumP + pow2) % MOD
       }
       pow2 = (pow2 * 2) % MOD
   }
   // pow2 now == 2^m mod MOD
   twoPowM := pow2

   // compute geometric series GS = sum_{j=0..k-1} (2^m)^j mod MOD
   var GS int64
   if twoPowM == 1 {
       // denominator zero, series is k mod MOD
       GS = k % MOD
   } else {
       numerator := (modPow(twoPowM, k) - 1 + MOD) % MOD
       denom := (twoPowM - 1 + MOD) % MOD
       GS = numerator * modInv(denom) % MOD
   }
   ans := GS * sumP % MOD
   // output result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, ans)
}
