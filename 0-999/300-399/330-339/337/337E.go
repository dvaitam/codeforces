package main

import (
   "bufio"
   "fmt"
   "math/big"
   "math/rand"
   "os"
   "time"
)

// fast factorization for 64-bit integers
// mulMod computes (a * b) % mod safely
func mulMod(a, b, mod uint64) uint64 {
   res := new(big.Int).Mul(new(big.Int).SetUint64(a), new(big.Int).SetUint64(b))
   res.Mod(res, new(big.Int).SetUint64(mod))
   return res.Uint64()
}
func powMod(a, d, mod uint64) uint64 {
   res := uint64(1)
   for d > 0 {
       if d&1 == 1 {
           res = mulMod(res, a, mod)
       }
       a = mulMod(a, a, mod)
       d >>= 1
   }
   return res
}
func isPrime(n uint64) bool {
   if n < 2 {
       return false
   }
   // small primes
   small := []uint64{2,3,5,7,11,13,17,19,23,29,31,37}
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
   // bases for deterministic Miller-Rabin up to 2^64
   bases := []uint64{2, 325, 9375, 28178, 450775, 9780504, 1795265022}
   for _, a := range bases {
       if a%n == 0 {
           continue
       }
       x := powMod(a, d, n)
       if x == 1 || x == n-1 {
           continue
       }
       skip := false
       for r := 1; r < s; r++ {
           x = mulMod(x, x, n)
           if x == n-1 {
               skip = true
               break
           }
       }
       if skip {
           continue
       }
       return false
   }
   return true
}
// Pollard's Rho
func pollardsRho(n uint64) uint64 {
   if n%2 == 0 {
       return 2
   }
   if n%3 == 0 {
       return 3
   }
   for {
       c := uint64(rand.Int63n(int64(n-1))) + 1
       x := uint64(rand.Int63n(int64(n)))
       y := x
       d := uint64(1)
       for d == 1 {
           x = (mulMod(x, x, n) + c) % n
           y = (mulMod(y, y, n) + c) % n
           y = (mulMod(y, y, n) + c) % n
           if x > y {
               d = gcd(x-y, n)
           } else {
               d = gcd(y-x, n)
           }
           if d == n {
               break
           }
       }
       if d > 1 && d < n {
           return d
       }
   }
}
func gcd(a, b uint64) uint64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}
// factor returns prime factors (with multiplicity)
func factor(n uint64, res *[]uint64) {
   if n == 1 {
       return
   }
   if isPrime(n) {
       *res = append(*res, n)
       return
   }
   d := pollardsRho(n)
   factor(d, res)
   factor(n/d, res)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   a := make([]uint64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   rand.Seed(time.Now().UnixNano())
   // compute Omega for each
   omega := make([]int, n)
   for i := 0; i < n; i++ {
       var fs []uint64
       factor(a[i], &fs)
       omega[i] = len(fs)
   }
   // P_j: possible parents for j if some a[i] divisible by a[j]
   isChildable := make([]bool, n)
   for j := 0; j < n; j++ {
       for i := 0; i < n; i++ {
           if i != j && a[i]%a[j] == 0 {
               isChildable[j] = true
               break
           }
       }
   }
   // classify B (childable) and Un (non-childable)
   B := make([]int, 0, n)
   unCount := 0
   for i := 0; i < n; i++ {
       if isChildable[i] {
           B = append(B, i)
       } else {
           unCount++
       }
   }
   // sumOmegaComposite
   sumOmegaComp := 0
   for i := 0; i < n; i++ {
       if omega[i] > 1 {
           sumOmegaComp += omega[i]
       }
   }
   best := int(1e18)
   m := len(B)
   // iterate masks
   for mask := 0; mask < (1 << m); mask++ {
       // sum omega selected
       sumSel := 0
       bits := 0
       for k := 0; k < m; k++ {
           if mask&(1<<k) != 0 {
               sumSel += omega[B[k]]
               bits++
           }
       }
       roots := unCount + (m - bits)
       penalty := 0
       if roots > 1 {
           penalty = 1
       }
       total := n + sumOmegaComp - sumSel + penalty
       if total < best {
           best = total
       }
   }
   fmt.Println(best)
}
