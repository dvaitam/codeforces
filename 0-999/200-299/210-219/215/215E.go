package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var l, r uint64
   if _, err := fmt.Fscan(in, &l, &r); err != nil {
       return
   }
   fmt.Println(countPeriodic(r) - countPeriodic(l-1))
}

// countPeriodic returns number of periodic numbers in [1..x]
func countPeriodic(x uint64) uint64 {
   if x == 0 {
       return 0
   }
   // get bits of x (1-based index)
   var bits []int
   for b := 63; b >= 0; b-- {
       if x>>b > 0 {
           for i := b; i >= 0; i-- {
               bits = append(bits, int((x>>i)&1))
           }
           break
       }
   }
   N := len(bits)
   // precompute mu and A up to N
   mu := computeMu(N)
   A := make([]uint64, N+1)
   for n := 1; n <= N; n++ {
       // P(n) = sum_{d|n} mu[d] * 2^{n/d}
       var P int64
       for d := 1; d <= n; d++ {
           if n%d == 0 {
               P += int64(mu[d]) * int64(1<<(n/d))
           }
       }
       // A[n] = number of periodic strings length n (s[0]=1)
       // = 2^{n-1} - P/2
       A[n] = (uint64(1)<<(n-1) - uint64(P/2))
   }
   // sum for lengths < N
   var res uint64
   for n := 1; n < N; n++ {
       res += A[n]
   }
   // handle length == N
   // divisors of N excluding N
   var D []int
   for d := 1; d < N; d++ {
       if N%d == 0 {
           D = append(D, d)
       }
   }
   m := len(D)
   // inclusion-exclusion coefficients C[g]
   C := make(map[int]int)
   // mask over D
   for mask := 1; mask < (1 << m); mask++ {
       var g int
       bitsCount := 0
       for j := 0; j < m; j++ {
           if mask>>j&1 == 1 {
               bitsCount++
               if g == 0 {
                   g = D[j]
               } else {
                   // gcd
                   a, b := g, D[j]
                   for b != 0 {
                       a, b = b, a%b
                   }
                   g = a
               }
           }
       }
       sign := 1
       if bitsCount%2 == 0 {
           sign = -1
       }
       C[g] += sign
   }
   // count for each g
   for g, coef := range C {
       if coef == 0 {
           continue
       }
       cnt := countWithPeriod(bits, g)
       // coef may be negative
       res += uint64(coef) * cnt
   }
   return res
}

// computeMu computes Möbius function mu[1..n]
func computeMu(n int) []int {
   mu := make([]int, n+1)
   mu[1] = 1
   for i := 2; i <= n; i++ {
       mu[i] = 1
       x := i
       var p int
       for d := 2; d*d <= x; d++ {
           if x%d == 0 {
               p++
               x /= d
               if x%d == 0 {
                   mu[i] = 0
                   break
               }
           }
           for x%d == 0 {
               x /= d
           }
       }
       if mu[i] == 0 {
           continue
       }
       if x > 1 {
           p++
       }
       if p%2 == 1 {
           mu[i] = -1
       } else {
           mu[i] = 1
       }
   }
   return mu
}

// countWithPeriod counts t of length g (t[0]=1) such that s (repeated) <= bits
func countWithPeriod(bits []int, g int) uint64 {
   N := len(bits)
   // count equal case
   eq := true
   for i := 0; i < N; i++ {
       if bits[i] != bits[i%g] {
           eq = false
           break
       }
   }
   var cnt uint64
   if eq && bits[0] == 1 {
       cnt = 1
   }
   // count by first difference at i
   for i := 0; i < N; i++ {
       if bits[i] == 0 {
           continue
       }
       forced := make([]int8, g)
       for k := 0; k < g; k++ {
           forced[k] = -1
       }
       ok := true
       // apply constraints from j< i
       for j := 0; j < i; j++ {
           k := j % g
           v := int8(bits[j])
           if forced[k] == -1 {
               forced[k] = v
           } else if forced[k] != v {
               ok = false
               break
           }
       }
       if !ok {
           continue
       }
       ki := i % g
       // enforce t[ki] = 0
       if forced[ki] != -1 {
           if forced[ki] == 1 {
               continue
           }
       } else {
           forced[ki] = 0
       }
       // count free bits
       used := 0
       for k := 0; k < g; k++ {
           if forced[k] != -1 {
               used++
           }
       }
       free := g - used
       cnt += uint64(1) << free
   }
   return cnt
}
