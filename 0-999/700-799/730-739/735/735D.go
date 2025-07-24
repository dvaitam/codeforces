package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n uint64
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   switch {
   case isPrime(n):
       fmt.Println(1)
   case n%2 == 0:
       fmt.Println(2)
   case isPrime(n - 2):
       fmt.Println(2)
   default:
       fmt.Println(3)
   }
}

// isPrime checks if n is a prime using deterministic Miller-Rabin
func isPrime(n uint64) bool {
   if n < 2 {
       return false
   }
   smallPrimes := []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}
   for _, p := range smallPrimes {
       if n == p {
           return true
       }
       if n%p == 0 {
           return n == p
       }
   }
   d := n - 1
   s := 0
   for d&1 == 0 {
       d >>= 1
       s++
   }
   // bases for deterministic check for n < 2^32
   bases := []uint64{2, 7, 61}
   for _, a := range bases {
       if a >= n {
           continue
       }
       if !millerTest(a, d, s, n) {
           return false
       }
   }
   return true
}

// millerTest performs Miller-Rabin test for base a
func millerTest(a, d, s, n uint64) bool {
   x := modPow(a, d, n)
   if x == 1 || x == n-1 {
       return true
   }
   for i := 1; i < s; i++ {
       x = mulMod(x, x, n)
       if x == n-1 {
           return true
       }
   }
   return false
}

// modPow computes a^e mod m
func modPow(a, e, m uint64) uint64 {
   result := uint64(1)
   a %= m
   for e > 0 {
       if e&1 == 1 {
           result = mulMod(result, a, m)
       }
       a = mulMod(a, a, m)
       e >>= 1
   }
   return result
}

// mulMod computes (a * b) % mod safely
func mulMod(a, b, mod uint64) uint64 {
   return (a * b) % mod
}
