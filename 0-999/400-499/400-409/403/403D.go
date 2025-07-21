package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007
const MAXN = 1000
// Kmax such that k(k-1)/2 + k <= MAXN => k <= 44 for MAXN=1000
const Kmax = 44

func add(a, b int) int {
   a += b
   if a >= MOD {
       a -= MOD
   }
   return a
}

func mul(a, b int) int {
   return int((int64(a) * int64(b)) % MOD)
}

func main() {
   // precompute D[k][s]: number of ways to choose k distinct non-negative ints summing to s
   var D [Kmax+1][MAXN+1]int
   D[0][0] = 1
   for d := 0; d <= MAXN; d++ {
       for k := Kmax; k >= 1; k-- {
           for s := d; s <= MAXN; s++ {
               if D[k-1][s-d] != 0 {
                   D[k][s] = add(D[k][s], D[k-1][s-d])
               }
           }
       }
   }
   // factorials and inverse factorials for combinations
   fact := make([]int, MAXN+1)
   invfact := make([]int, MAXN+1)
   fact[0] = 1
   for i := 1; i <= MAXN; i++ {
       fact[i] = mul(fact[i-1], i)
   }
   invfact[MAXN] = modInv(fact[MAXN])
   for i := MAXN; i > 0; i-- {
       invfact[i-1] = mul(invfact[i], i)
   }
   // precompute C(n,k) on the fly via fact
   comb := func(n, k int) int {
       if n < 0 || k < 0 || n < k {
           return 0
       }
       return mul(fact[n], mul(invfact[k], invfact[n-k]))
   }
   // precompute k! for k<=Kmax
   factk := make([]int, Kmax+1)
   factk[0] = 1
   for i := 1; i <= Kmax; i++ {
       factk[i] = mul(factk[i-1], i)
   }
   // f[n][k]
   var f [MAXN+1][Kmax+1]int
   for k := 1; k <= Kmax; k++ {
       // minimal sum s is k(k-1)/2
       minS := k * (k - 1) / 2
       for s := minS; s <= MAXN; s++ {
           cnt := D[k][s]
           if cnt == 0 {
               continue
           }
           coef := mul(factk[k], cnt)
           // n must satisfy n-k >= s => n >= k+s
           for n := k + s; n <= MAXN; n++ {
               // choose gaps: C(n - s, k)
               f[n][k] = add(f[n][k], mul(coef, comb(n-s, k)))
           }
       }
   }
   // input
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var n, k int
       fmt.Fscan(reader, &n, &k)
       if k < 0 || k > Kmax || k > n {
           fmt.Fprintln(writer, 0)
       } else {
           fmt.Fprintln(writer, f[n][k])
       }
   }
}

// modInv computes modular inverse of a under MOD
func modInv(a int) int {
   return modPow(a, MOD-2)
}

func modPow(a, e int) int {
   res := 1
   base := a % MOD
   for e > 0 {
       if e&1 != 0 {
           res = mul(res, base)
       }
       base = mul(base, base)
       e >>= 1
   }
   return res
}
