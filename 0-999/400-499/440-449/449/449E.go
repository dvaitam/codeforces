package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007
const MAXN = 1000000

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   // Read all queries
   ns := make([]int, t)
   ms := make([]int, t)
   maxK := 0
   for i := 0; i < t; i++ {
       fmt.Fscan(reader, &ns[i], &ms[i])
       k := ns[i]
       if ms[i] < k {
           k = ms[i]
       }
       if k > maxK {
           maxK = k
       }
   }
   if maxK > MAXN {
       maxK = MAXN
   }
   // Precompute phi
   phi := make([]int, maxK+1)
   for i := 0; i <= maxK; i++ {
       phi[i] = i
   }
   for i := 2; i <= maxK; i++ {
       if phi[i] == i {
           for j := i; j <= maxK; j += i {
               phi[j] -= phi[j] / i
           }
       }
   }
   // Precompute sum of gcds f2[k] = sum_{i=1..k} gcd(i,k)
   f2 := make([]int64, maxK+1)
   for i := 1; i <= maxK; i++ {
       pi := int64(phi[i])
       for k0 := i; k0 <= maxK; k0 += i {
           // k0 = i * d => d = k0 / i
           d := int64(k0 / i)
           f2[k0] += d * pi
       }
   }
   // Precompute G[k], and prefix sums A,B,C
   G := make([]int64, maxK+1)
   A := make([]int64, maxK+1)
   B := make([]int64, maxK+1)
   C := make([]int64, maxK+1)
   inv2 := int64((MOD + 1) / 2)
   inv6 := int64(166666668) // inverse of 6 mod MOD
   for k := 1; k <= maxK; k++ {
       kk := int64(k)
       // T1 = k*(k+1)*(2k+1)/6
       t1 := kk * (kk + 1) % MOD * (2*kk + 1) % MOD * inv6 % MOD
       // T3 = k*(k+1)/2
       t3 := kk * (kk + 1) % MOD * inv2 % MOD
       sg := (f2[k] + kk) % MOD
       // G[k] = sum_{dx=0..k} f(dx,k-dx); we need only dx>=1, so subtract f(0,k)=k^2
       gk := (2*t1%MOD - 4*t3%MOD + 2*sg%MOD) % MOD
       // subtract f(0,k) = k^2 to exclude dx=0 case (duplicate axis-aligned)
       gk = (gk - (kk*kk%MOD) + MOD) % MOD
       G[k] = gk
       A[k] = (A[k-1] + gk) % MOD
       B[k] = (B[k-1] + gk*kk) % MOD
       C[k] = (C[k-1] + gk*kk%MOD*kk%MOD) % MOD
   }
   // Answer queries
   for i := 0; i < t; i++ {
       n, m := ns[i], ms[i]
       K := n
       if m < K {
           K = m
       }
       if K <= 0 {
           fmt.Fprintln(writer, 0)
           continue
       }
       np := int64(n + 1)
       mp := int64(m + 1)
       s0 := A[K]
       s1 := B[K]
       s2 := C[K]
       // ans = (n+1)(m+1)*s0 - (n+m+2)*s1 + s2
       ans := (np * mp % MOD * s0 % MOD - (np+mp)%MOD * s1 % MOD + s2) % MOD
       if ans < 0 {
           ans += MOD
       }
       fmt.Fprintln(writer, ans)
   }
}
