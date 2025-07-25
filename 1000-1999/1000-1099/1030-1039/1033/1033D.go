package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "math/rand"
   "os"
   "time"
)

const mod = 998244353

// mulMod computes (a * b) % m safely for uint64
func mulMod(a, b, m uint64) uint64 {
   hi, lo := bits.Mul64(a, b)
   if hi == 0 {
       return lo % m
   }
   // compute 2^64 % m
   t := (uint64(1) << 32) % m
   two64 := (t * t) % m
   return ((hi % m) * two64 % m + lo % m) % m
}

// powMod computes a^e % m
func powMod(a, e, m uint64) uint64 {
   res := uint64(1)
   a %= m
   for e > 0 {
       if e&1 == 1 {
           res = mulMod(res, a, m)
       }
       a = mulMod(a, a, m)
       e >>= 1
   }
   return res
}

// isPrime tests n for primality (deterministic for 64-bit)
func isPrime(n uint64) bool {
   if n < 2 {
       return false
   }
   // small primes
   small := []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
   for _, p := range small {
       if n%p == 0 {
           return n == p
       }
   }
   // write n-1 = d * 2^s
   d := n - 1
   s := 0
   for d&1 == 0 {
       d >>= 1
       s++
   }
   bases := []uint64{2, 325, 9375, 28178, 450775, 9780504, 1795265022}
   for _, a := range bases {
       if a%n == 0 {
           continue
       }
       x := powMod(a, d, n)
       if x == 1 || x == n-1 {
           continue
       }
       composite := true
       for r := 1; r < s; r++ {
           x = mulMod(x, x, n)
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

// pollardsRho returns a non-trivial factor of n
func pollardsRho(n uint64) uint64 {
   if n%2 == 0 {
       return 2
   }
   if isPrime(n) {
       return n
   }
   randSrc := rand.New(rand.NewSource(time.Now().UnixNano()))
   for {
       c := randSrc.Uint64()%(n-1) + 1
       x := randSrc.Uint64()%n
       y := x
       d := uint64(1)
       for d == 1 {
           x = (mulMod(x, x, n) + c) % n
           y = (mulMod(y, y, n) + c) % n
           y = (mulMod(y, y, n) + c) % n
           // |x-y|
           var diff uint64
           if x > y {
               diff = x - y
           } else {
               diff = y - x
           }
           d = gcd(diff, n)
           if d == n {
               break
           }
       }
       if d > 1 && d < n {
           return d
       }
   }
}

// gcd returns greatest common divisor of a and b
func gcd(a, b uint64) uint64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

// factorize recursively factors n into primes, populating mp
func factorize(n uint64, mp map[uint64]uint64) {
   if n == 1 {
       return
   }
   if isPrime(n) {
       mp[n]++
       return
   }
   d := pollardsRho(n)
   factorize(d, mp)
   factorize(n/d, mp)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   total := make(map[uint64]uint64)
   for i := 0; i < n; i++ {
       var a uint64
       fmt.Fscan(in, &a)
       mp := make(map[uint64]uint64)
       factorize(a, mp)
       for p, cnt := range mp {
           total[p] += cnt
       }
   }
   ans := uint64(1)
   for _, cnt := range total {
       ans = ans * (cnt + 1) % mod
   }
   fmt.Println(ans)
}
