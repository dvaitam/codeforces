package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

// modpow computes a^e mod mod
func modpow(a, e int64) int64 {
   res := int64(1)
   a %= mod
   for e > 0 {
       if e&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       e >>= 1
   }
   return res
}

// modinv computes modular inverse of a mod mod
func modinv(a int64) int64 {
   return modpow(a, mod-2)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   l := make([]int64, n)
   r := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &l[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &r[i])
   }
   // lengths and inverses
   lenv := make([]int64, n)
   inv := make([]int64, n)
   for i := 0; i < n; i++ {
       lenv[i] = r[i] - l[i] + 1
       inv[i] = modinv(lenv[i] % mod)
   }
   // probabilities q[i] = P(x[i]==x[i-1]) for i=1..n-1
   q := make([]int64, n)
   for i := 1; i < n; i++ {
       // intersection count
       low := l[i]
       if l[i-1] > low {
           low = l[i-1]
       }
       high := r[i]
       if r[i-1] < high {
           high = r[i-1]
       }
       cnt := int64(0)
       if high >= low {
           cnt = high - low + 1
       }
       // q
       tmp := cnt % mod
       tmp = tmp * inv[i-1] % mod
       tmp = tmp * inv[i] % mod
       q[i] = tmp
   }
   // accumulate sums
   var sumP, sumP2, sumSkip, sumAdjE int64
   // p[i] = 1 - q[i]
   prevP := int64(0)
   // compute sumP and sumP2
   for i := 1; i < n; i++ {
       pi := (1 - q[i] + mod) % mod
       sumP = (sumP + pi) % mod
       sumP2 = (sumP2 + pi*pi%mod) % mod
       if i > 1 {
           // skip term p[i-1]*p[i]
           sumSkip = (sumSkip + prevP*pi) % mod
           // adjacent E term at k=i-1: between q[i-1], q[i]
           // compute t for triple intersection at positions i-2, i-1, i
           low := l[i]
           if l[i-2] > low {
               low = l[i-2]
           }
           if l[i-1] > low {
               low = l[i-1]
           }
           high := r[i]
           if r[i-2] < high {
               high = r[i-2]
           }
           if r[i-1] < high {
               high = r[i-1]
           }
           tcnt := int64(0)
           if high >= low {
               tcnt = high - low + 1
           }
           t := tcnt % mod
           t = t * inv[i-2] % mod
           t = t * inv[i-1] % mod
           t = t * inv[i] % mod
           adj := (1 - q[i-1] - q[i] + t) % mod
           if adj < 0 {
               adj += mod
           }
           sumAdjE = (sumAdjE + adj) % mod
       }
       prevP = pi
   }
   // compute result: E(B^2) = 1 + 3*sumP + sumP^2 - sumP2 + 2*sumAdjE - 2*sumSkip
   res := (1 + 3*sumP%mod + sumP*sumP%mod - sumP2 + 2*sumAdjE%mod - 2*sumSkip%mod) % mod
   if res < 0 {
       res += mod
   }
   // output
   fmt.Println(res)
}
