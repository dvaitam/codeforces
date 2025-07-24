package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, d int
   var mod int64
   fmt.Fscan(in, &n, &d, &mod)
   MOD := mod
   // Precompute inv factorials up to max d
   maxC := d
   invFact := make([]int64, maxC+1)
   fact := make([]int64, maxC+1)
   fact[0] = 1
   for i := 1; i <= maxC; i++ {
       fact[i] = fact[i-1] * int64(i) % MOD
   }
   invFact[maxC] = modInv(fact[maxC], MOD)
   for i := maxC; i > 0; i-- {
       invFact[i-1] = invFact[i] * int64(i) % MOD
   }
   inv2 := (MOD + 1) / 2
   // Prepare arrays H and R
   N := n
   H := make([]int64, N+1)
   R := make([]int64, N+1)
   // dpH[k][s]: using types processed, select k children sum sizes s, for c=d-1
   cH := d - 1
   dpH := make([][]int64, cH+1)
   for i := range dpH {
       dpH[i] = make([]int64, N+1)
   }
   // dpR for c=d
   cR := d
   dpR := make([][]int64, cR+1)
   for i := range dpR {
       dpR[i] = make([]int64, N+1)
   }
   // initialize with type size 1, H[1]=1
   H[1] = 1
   // process type m=1 in dpH and dpR
   dpH[0][0] = 1
   dpR[0][0] = 1
   // include m=1
   for t := 1; t <= cH; t++ {
       // comb( H[1]+t-1, t ) = comb(1+t-1,t)=1
       for k := cH - t; k >= 0; k-- {
           for s := 0; s + t <= N; s++ {
               dpH[k+t][s+t] = (dpH[k+t][s+t] + dpH[k][s]) % MOD
           }
       }
   }
   for t := 1; t <= cR; t++ {
       for k := cR - t; k >= 0; k-- {
           for s := 0; s + t <= N; s++ {
               dpR[k+t][s+t] = (dpR[k+t][s+t] + dpR[k][s]) % MOD
           }
       }
   }
   // R[1] = H[1] + dpR[d][0]
   R[1] = (H[1] + dpR[cR][0]) % MOD
   // Build for m=2..N
   for m := 2; m <= N; m++ {
       // compute H[m] = dpH[cH][m-1]
       if m-1 >= 0 {
           H[m] = dpH[cH][m-1]
       }
       // compute R[m]
       var cntR int64
       if m-1 >= 0 {
           cntR = dpR[cR][m-1]
       }
       R[m] = (H[m] + cntR) % MOD
       // prepare ways for type m
       // for dpH, t up to cH
       waysH := make([]int64, cH+1)
       for t := 1; t <= cH; t++ {
           // comb(H[m]+t-1, t)
           waysH[t] = combSeq(H[m], t, invFact, MOD)
       }
       // for dpR, t up to cR
       waysR := make([]int64, cR+1)
       for t := 1; t <= cR; t++ {
           waysR[t] = combSeq(H[m], t, invFact, MOD)
       }
       // update dpH with type m
       for t := 1; t <= cH; t++ {
           w := waysH[t]
           if w == 0 {
               continue
           }
           for k := cH - t; k >= 0; k-- {
               for s := 0; s + m*t <= N; s++ {
                   dpH[k+t][s+m*t] = (dpH[k+t][s+m*t] + dpH[k][s]*w) % MOD
               }
           }
       }
       // update dpR with type m
       for t := 1; t <= cR; t++ {
           w := waysR[t]
           if w == 0 {
               continue
           }
           for k := cR - t; k >= 0; k-- {
               for s := 0; s + m*t <= N; s++ {
                   dpR[k+t][s+m*t] = (dpR[k+t][s+m*t] + dpR[k][s]*w) % MOD
               }
           }
       }
   }
   // compute convolution sum
   var conv int64
   for i := 1; i < n; i++ {
       conv = (conv + R[i]*R[n-i]) % MOD
   }
   var mid int64
   if n%2 == 0 {
       mid = R[n/2]
   }
   // U[n] = R[n] - (conv + mid) * inv2
   ans := (R[n] - (conv+mid)%MOD*inv2%MOD) % MOD
   if ans < 0 {
       ans += MOD
   }
   fmt.Println(ans)
}

// modInv computes modular inverse of a mod m (m prime)
func modInv(a, m int64) int64 {
   return modPow(a, m-2, m)
}

// modPow computes a^e mod m
func modPow(a, e, m int64) int64 {
   res := int64(1)
   a %= m
   for e > 0 {
       if e&1 == 1 {
           res = res * a % m
       }
       a = a * a % m
       e >>= 1
   }
   return res
}

// combSeq computes C(x+t-1, t) = (x)*(x+1)*...*(x+t-1)/t!
func combSeq(x int64, t int, invFact []int64, MOD int64) int64 {
   // numerator = x*(x+1)*...*(x+t-1)
   num := int64(1)
   for i := 0; i < t; i++ {
       num = num * ((x + int64(i)) % MOD) % MOD
   }
   return num * invFact[t] % MOD
}
