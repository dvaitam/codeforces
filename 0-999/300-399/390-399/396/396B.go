package main

import (
   "bufio"
   "fmt"
   "os"
)

// modMul returns (a*b) % mod, safe for mod < 2^32
func modMul(a, b, mod uint64) uint64 {
   return (a * b) % mod
}

// modPow returns (a^d) % mod
func modPow(a, d, mod uint64) uint64 {
   res := uint64(1)
   a %= mod
   for d > 0 {
       if d&1 == 1 {
           res = modMul(res, a, mod)
       }
       a = modMul(a, a, mod)
       d >>= 1
   }
   return res
}

// isPrime tests primality for n < 2^32 using deterministic Miller-Rabin
func isPrime(n uint64) bool {
   if n < 2 {
       return false
   }
   // small primes
   smallPrimes := []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
   for _, p := range smallPrimes {
       if n == p {
           return true
       }
       if n%p == 0 {
           return n == p
       }
   }
   // write n-1 as d*2^s
   d := n - 1
   s := 0
   for d&1 == 0 {
       d >>= 1
       s++
   }
   // test bases
   bases := []uint64{2, 7, 61}
   for _, a := range bases {
       if a >= n {
           continue
       }
       x := modPow(a, d, n)
       if x == 1 || x == n-1 {
           continue
       }
       composite := true
       for r := 1; r < s; r++ {
           x = modMul(x, x, n)
           if x == n-1 {
               composite = false
               break
           }
       }
       if composite {
           return false
       }
   }
   return true
}

// gcd returns greatest common divisor of a and b
func gcd(a, b uint64) uint64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var n uint64
       fmt.Fscan(reader, &n)
       // find p_k = largest prime <= n
       pk := n
       for pk >= 2 && !isPrime(pk) {
           pk--
       }
       // find p_{k+1} = smallest prime > n
       pk1 := n + 1
       for !isPrime(pk1) {
           pk1++
       }
       // compute R = n - pk + 1 - pk1
       R := int64(n - pk + 1) - int64(pk1)
       // D = pk * pk1
       D := pk * pk1
       // numerator p = D + 2*R
       // since R may be negative, compute as int64 then convert
       pnum := int64(D) + 2*R
       // denominator q = 2*D
       qden := uint64(2) * D
       // make positive numerator
       qnum := uint64(pnum)
       g := gcd(qnum, qden)
       qnum /= g
       qden /= g
       fmt.Fprintf(writer, "%d/%d\n", qnum, qden)
   }
}
