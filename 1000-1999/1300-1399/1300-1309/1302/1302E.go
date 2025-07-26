package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD int64 = 1234567891

// extGCD returns gcd(a,b) and x,y such that ax + by = gcd(a,b)
func extGCD(a, b int64) (int64, int64, int64) {
   if b == 0 {
       return a, 1, 0
   }
   g, x1, y1 := extGCD(b, a%b)
   x := y1
   y := x1 - (a/b)*y1
   return g, x, y
}

// modInv returns modular inverse of a mod MOD; assumes gcd(a,MOD)=1
func modInv(a int64) int64 {
   _, x, _ := extGCD(a, MOD)
   inv := x % MOD
   if inv < 0 {
       inv += MOD
   }
   return inv
}

// modPow computes a^e mod MOD
func modPow(a, e int64) int64 {
   res := int64(1)
   a %= MOD
   for e > 0 {
       if e&1 == 1 {
           res = (res * a) % MOD
       }
       a = (a * a) % MOD
       e >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, a, b int64
   if _, err := fmt.Fscan(in, &n, &a, &b); err != nil {
       return
   }
   // Probability p = a/b, q = (b-a)/b
   invB := modInv(b)
   p := (a % MOD) * invB % MOD
   q := ((b-a)%MOD+MOD)%MOD * invB % MOD
   // Compute p^n and q^n
   pN := modPow(p, n)
   qN := modPow(q, n)
   // Answer = 1 - p^n - q^n mod MOD
   ans := (1 - pN - qN) % MOD
   if ans < 0 {
       ans += MOD
   }
   fmt.Println(ans)
}
