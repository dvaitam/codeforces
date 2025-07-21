package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func modexp(a, e int64) int64 {
   res := int64(1)
   a %= MOD
   for e > 0 {
       if e&1 == 1 {
           res = res * a % MOD
       }
       a = a * a % MOD
       e >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   var p int64
   fmt.Fscan(reader, &n, &p)
   // compute Grundy up to N
   const s0 = 1000
   const Tmax = 2000
   N := s0 + 2*Tmax
   g := make([]int, N+1)
   for L := 0; L <= 2; L++ {
       if L <= N {
           g[L] = 0
       }
   }
   for L := 3; L <= N; L++ {
       d := L / 3
       a := g[d]
       b := g[L-d]
       // mex of {a, b}
       if a != 0 && b != 0 {
           g[L] = 0
       } else if a != 1 && b != 1 {
           g[L] = 1
       } else {
           g[L] = 2
       }
   }
   // detect period T for g[L] for L>=s0
   var T int
   for t := 1; t <= Tmax; t++ {
       ok := true
       for i := 0; i < t; i++ {
           if s0+2*t > N {
               ok = false
               break
           }
           if g[s0+i] != g[s0+t+i] || g[s0+i] != g[s0+2*t+i] {
               ok = false
               break
           }
       }
       if ok {
           T = t
           break
       }
   }
   if T == 0 {
       T = N // no period
   }
   // compute C_k
   C := make([]int64, 3)
   // small L < s0 or if p-1 < s0
   maxL := p - 1
   if maxL < 0 {
       maxL = 0
   }
   for L := int64(1); L <= maxL && L < s0; L++ {
       k := g[L]
       C[k] = (C[k] + (p - L) % MOD) % MOD
   }
   if maxL >= s0 {
       total := maxL - s0 + 1 // = p-s0
       M := total / int64(T)
       rem := int(total % int64(T))
       // precompute per i in [0,T)
       // count per k
       cnt := make([]int64, 3)
       sum_si := make([]int64, 3)
       for i := 0; i < T; i++ {
           L := int64(s0 + i)
           k := g[s0+i]
           cnt[k]++
           sum_si[k] = (sum_si[k] + L) % MOD
       }
       inv2 := int64((MOD + 1) / 2)
       // sum over full periods
       for k := 0; k <= 2; k++ {
           // M * sum over j: (p - (s0+i) - j*T) for each i
           // sum_full = M*(cnt[k]*p - sum_si[k]) - T*(M*(M-1)/2)%MOD * cnt[k]
           term1 := (cnt[k]* (p%MOD) % MOD - sum_si[k] + MOD) % MOD
           sum_full := (term1 * (M % MOD)) % MOD
           mMn1 := (M % MOD) * ((M - 1) % MOD) % MOD * inv2 % MOD
           sum_full = (sum_full - (int64(T)%MOD)*mMn1%MOD*cnt[k]%MOD + MOD) % MOD
           C[k] = (C[k] + sum_full) % MOD
       }
       // sum rem
       for i := 0; i < rem; i++ {
           L := int64(s0 + i + int(M)*T)
           k := g[s0+i]
           C[k] = (C[k] + (p - L) % MOD) % MOD
       }
   }
   // total_all = (C0+C1+C2)^n
   sumC := (C[0] + C[1] + C[2]) % MOD
   total_all := modexp(sumC, int64(n))
   // XOR convolution to get zero XOR count
   size := 4
   A := make([]int64, size)
   A[0] = C[0]
   A[1] = C[1]
   A[2] = C[2]
   A[3] = 0
   // FWHT
   for len := 1; len < size; len <<= 1 {
       for i := 0; i < size; i += len * 2 {
           for j := 0; j < len; j++ {
               u := A[i+j]
               v := A[i+j+len]
               A[i+j] = (u + v) % MOD
               A[i+j+len] = (u - v + MOD) % MOD
           }
       }
   }
   for i := 0; i < size; i++ {
       A[i] = modexp(A[i], int64(n))
   }
   // inverse FWHT
   for len := 1; len < size; len <<= 1 {
       for i := 0; i < size; i += len * 2 {
           for j := 0; j < len; j++ {
               u := A[i+j]
               v := A[i+j+len]
               A[i+j] = (u + v) * int64((MOD+1)/2) % MOD
               A[i+j+len] = (u - v + MOD) * int64((MOD+1)/2) % MOD
           }
       }
   }
   zero := A[0]
   ans := (total_all - zero + MOD) % MOD
   fmt.Fprintln(writer, ans)
}
