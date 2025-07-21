package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // sieve primes up to sqrt(max ai) ~31623
   maxP := 31623
   isPrime := make([]bool, maxP+1)
   for i := 2; i <= maxP; i++ {
       isPrime[i] = true
   }
   primes := make([]int, 0, 3400)
   for i := 2; i*i <= maxP; i++ {
       if isPrime[i] {
           for j := i * i; j <= maxP; j += i {
               isPrime[j] = false
           }
       }
   }
   for i := 2; i <= maxP; i++ {
       if isPrime[i] {
           primes = append(primes, i)
       }
   }
   // factorize and sum exponents
   exp := make(map[int64]int)
   var maxE int
   for _, v := range a {
       x := v
       for _, p := range primes {
           pp := int64(p)
           if pp*pp > x {
               break
           }
           if x%pp == 0 {
               cnt := 0
               for x%pp == 0 {
                   x /= pp
                   cnt++
               }
               exp[pp] += cnt
               if exp[pp] > maxE {
                   maxE = exp[pp]
               }
           }
       }
       if x > 1 {
           // x is prime
           exp[x]++
           if exp[x] > maxE {
               maxE = exp[x]
           }
       }
   }
   // prepare factorials up to maxE + n - 1
   lim := maxE + n - 1
   fact := make([]int64, lim+1)
   invf := make([]int64, lim+1)
   fact[0] = 1
   for i := 1; i <= lim; i++ {
       fact[i] = fact[i-1] * int64(i) % mod
   }
   invf[lim] = modInv(fact[lim])
   for i := lim; i > 0; i-- {
       invf[i-1] = invf[i] * int64(i) % mod
   }
   // compute result
   res := int64(1)
   for _, e := range exp {
       // C(e + n - 1, n - 1)
       top := e + n - 1
       ways := fact[top] * invf[n-1] % mod * invf[e] % mod
       res = res * ways % mod
   }
   fmt.Println(res)
}

// modInv computes modular inverse via Fermat's little theorem
func modInv(x int64) int64 {
   return modPow(x, mod-2)
}

// modPow computes x^e mod mod
func modPow(x, e int64) int64 {
   res := int64(1)
   x %= mod
   for e > 0 {
       if e&1 == 1 {
           res = res * x % mod
       }
       x = x * x % mod
       e >>= 1
   }
   return res
}
