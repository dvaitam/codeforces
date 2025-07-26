package main

import (
   "bufio"
   "fmt"
   "math"
   "math/rand"
   "os"
   "time"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   a := make([]uint64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // check initial gcd
   g := a[0]
   for i := 1; i < n; i++ {
       g = gcd(g, a[i])
       if g == 1 {
           break
       }
   }
   if g > 1 {
       fmt.Println(0)
       return
   }
   // random sampling of indices
   rand.Seed(time.Now().UnixNano())
   idx := rand.Perm(n)
   limit := 15
   if n < limit {
       limit = n
   }
   candPrimes := make(map[uint64]struct{})
   for i := 0; i < limit; i++ {
       for _, delta := range []int{-1, 0, 1} {
           v := int64(a[idx[i]]) + int64(delta)
           if v < 2 {
               continue
           }
           ps := factor(uint64(v))
           for _, p := range ps {
               candPrimes[p] = struct{}{}
           }
       }
   }
   best := uint64(math.MaxUint64)
   // evaluate each candidate prime
   for p := range candPrimes {
       if p <= 1 {
           continue
       }
       var cost uint64
       for i := 0; i < n; i++ {
           ai := a[i]
           r := ai % p
           if r != 0 {
               var c uint64
               if ai < p {
                   c = p - r
               } else {
                   // can subtract or add
                   dr := r
                   dadd := p - r
                   if dr < dadd {
                       c = dr
                   } else {
                       c = dadd
                   }
               }
               cost += c
               if cost >= best {
                   break
               }
           }
       }
       if cost < best {
           best = cost
       }
   }
   fmt.Println(best)
}

func gcd(a, b uint64) uint64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

// Miller-Rabin primality test
func isPrime(n uint64) bool {
   if n < 2 {
       return false
   }
   small := []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
   for _, p := range small {
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
   // bases for <2^64
   bases := []uint64{2, 325, 9375, 28178, 450775, 9780504, 1795265022}
   for _, a := range bases {
       if a%n == 0 {
           continue
       }
       x := modPow(a, d, n)
       if x == 1 || x == n-1 {
           continue
       }
       comp := true
       for r := 1; r < s; r++ {
           x = modMul(x, x, n)
           if x == n-1 {
               comp = false
               break
           }
       }
       if comp {
           return false
       }
   }
   return true
}

func modMul(a, b, mod uint64) uint64 {
   return (a * b) % mod
}

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

// Pollard's Rho algorithm
func pollardsRho(n uint64) uint64 {
   if n%2 == 0 {
       return 2
   }
   x := uint64(rand.Int63n(int64(n-2))) + 2
   y := x
   c := uint64(rand.Int63n(int64(n-1))) + 1
   d := uint64(1)
   for d == 1 {
       x = (modMul(x, x, n) + c) % n
       y = (modMul(y, y, n) + c) % n
       y = (modMul(y, y, n) + c) % n
       d = gcd((x>y?x-y:y-x), n)
       if d == n {
           return pollardsRho(n)
       }
   }
   return d
}

// factor returns prime factors (unique) of n
func factor(n uint64) []uint64 {
   var res []uint64
   var dfs func(uint64)
   dfs = func(n uint64) {
       if n == 1 {
           return
       }
       if isPrime(n) {
           res = append(res, n)
       } else {
           d := pollardsRho(n)
           dfs(d)
           dfs(n / d)
       }
   }
   dfs(n)
   // unique
   m := make(map[uint64]struct{})
   var uni []uint64
   for _, p := range res {
       if _, ok := m[p]; !ok {
           m[p] = struct{}{}
           uni = append(uni, p)
       }
   }
   return uni
}
