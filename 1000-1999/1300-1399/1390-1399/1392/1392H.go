package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353

// modexp computes a^e mod mod
func modexp(a, e int) int {
   res := 1
   base := a % mod
   for e > 0 {
       if e&1 == 1 {
           res = int((int64(res) * base) % mod)
       }
       base = int((int64(base) * base) % mod)
       e >>= 1
   }
   return res
}

// modinv computes modular inverse of a mod mod
func modinv(a int) int {
   // Fermat's little theorem, mod is prime
   return modexp((a%mod+mod)%mod, mod-2)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   // Precompute factorials and inv factorials
   fact := make([]int, n+1)
   invfact := make([]int, n+1)
   fact[0] = 1
   for i := 1; i <= n; i++ {
       fact[i] = int(int64(fact[i-1]) * i % mod)
   }
   invfact[n] = modinv(fact[n])
   for i := n; i > 0; i-- {
       invfact[i-1] = int(int64(invfact[i]) * i % mod)
   }
   // Precompute powers of (m+1) and m
   powA := make([]int, n+1)
   powB := make([]int, n+1)
   A := (m + 1) % mod
   B := m % mod
   powA[0], powB[0] = 1, 1
   for i := 1; i <= n; i++ {
       powA[i] = int(int64(powA[i-1]) * A % mod)
       powB[i] = int(int64(powB[i-1]) * B % mod)
   }
   // Sum S = sum_{j=1..n} (-1)^(j-1) * C(n,j) * powA[j] / (powA[j] - powB[j])
   var S int
   for j := 1; j <= n; j++ {
       // binomial C(n,j)
       bin := int(int64(fact[n]) * invfact[j] % mod * int64(invfact[n-j]) % mod)
       aj := powA[j]
       bj := powB[j]
       denom := aj - bj
       if denom < 0 {
           denom += mod
       }
       invDen := modinv(denom)
       term := int(int64(bin) * int64(aj) % mod * int64(invDen) % mod)
       if j%2 == 1 {
           S += term
           if S >= mod {
               S -= mod
           }
       } else {
           S -= term
           if S < 0 {
               S += mod
           }
       }
   }
   // Multiply by E[L] = (n + m + 1) / (m + 1)
   coeff := int(int64(n + m + 1) % mod * int64(modinv(m+1)) % mod)
   ans := int(int64(S) * int64(coeff) % mod)
   fmt.Println(ans)
}
