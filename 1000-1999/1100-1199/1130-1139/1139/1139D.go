package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func modPow(a, e int64) int64 {
   res := int64(1)
   a %= mod
   for e > 0 {
       if e&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       e >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var m int
   if _, err := fmt.Fscan(in, &m); err != nil {
       return
   }
   // Compute Möbius function with sieve
   mu := make([]int, m+1)
   p := make([]bool, m+1)
   primes := make([]int, m+1)
   mu[1] = 1
   pris := 0
   for i := 2; i <= m; i++ {
       if !p[i] {
           pris++
           primes[pris] = i
           mu[i] = 1
       }
       for j := 1; j <= pris; j++ {
           v := primes[j]
           s := i * v
           if s > m {
               break
           }
           p[s] = true
           if i%v != 0 {
               mu[s] = -mu[i]
           } else {
               mu[s] = 0
               break
           }
       }
   }
   var ansA, ansB int64 = 1, 1
   // accumulate contributions for i = 2..m (skip i=1 as per original logic)
   for i := 2; i <= m; i++ {
       if mu[i] == 0 {
           continue
       }
       // count multiples and non-multiples of i
       cntA := int64(m / i)
       cntB := int64(m - m/i)
       if mu[i] == -1 {
           cntA = (mod - cntA) % mod
       }
       // update numerator and denominator
       ansA = (ansA*cntB + ansB*cntA) % mod
       ansB = ansB * cntB % mod
   }
   // result = ansA / ansB mod
   res := ansA * modPow(ansB, mod-2) % mod
   fmt.Println(res)
}
