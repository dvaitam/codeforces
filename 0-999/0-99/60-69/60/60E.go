package main

import (
   "bufio"
   "fmt"
   "os"
)

// fastPowMod computes a^e mod m
func fastPowMod(a, e, m int64) int64 {
   a %= m
   var res int64 = 1
   for e > 0 {
       if e&1 == 1 {
           res = (res * a) % m
       }
       a = (a * a) % m
       e >>= 1
   }
   return res
}

// fibPair returns (F(n) mod m, F(n+1) mod m) using fast doubling
func fibPair(n, m int64) (int64, int64) {
   if n == 0 {
       return 0, 1
   }
   a, b := fibPair(n >> 1, m)
   // F(2k) = F(k) * [2*F(k+1) − F(k)]
   c := (a * ((2*b % m - a + m) % m)) % m
   // F(2k+1) = F(k)^2 + F(k+1)^2
   d := (a*a % m + b*b % m) % m
   if n&1 == 0 {
       return c, d
   }
   // F(2k+1), F(2k+2)
   return d, (c + d) % m
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var x, y, p int64
   if _, err := fmt.Fscan(in, &n, &x, &y, &p); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // special case n==1: no growth
   if n == 1 {
       fmt.Println(a[0] % p)
       return
   }
   // Precompute sums
   var midSumModP int64 = 0
   for i := 1; i+1 < n; i++ {
       midSumModP = (midSumModP + a[i]) % p
   }
   a0 := a[0]
   alast := a[n-1]
   // Phase1: compute S_x and M_x
   mod2p := p * 2
   // compute pow3_x mod 2p
   pow3x := fastPowMod(3, x, mod2p)
   // numerator for S_x: 2*midSum*3^x + (a0+alast)*(3^x+1)
   num := ( (2*midSumModP%mod2p)*pow3x % mod2p + ((a0+alast)%mod2p)*((pow3x+1)%mod2p) % mod2p ) % mod2p
   // S_x = num/2 mod p
   Sx := (num/2) % p
   // compute M_x: segment max for last pair
   // fib mod 2p
   Fx, Fx1 := fibPair(x, mod2p)
   L := a[n-2] % mod2p
   R := a[n-1] % mod2p
   Mx := (Fx*L + Fx1*R) % mod2p
   // Phase2: compute final S
   // D = a0 + Mx
   D := (a0 % mod2p + Mx) % mod2p
   // pow3_y mod 2p
   pow3y := fastPowMod(3, y, mod2p)
   // compute 2*S_final mod (2p): = pow3y*(2*Sx) - D*(pow3y-1)
   twoSx := (2 * Sx) % mod2p
   term1 := (pow3y * twoSx) % mod2p
   term2 := (D * ((pow3y - 1 + mod2p) % mod2p)) % mod2p
   twoSy := (term1 - term2 + mod2p) % mod2p
   // final S = twoSy/2 mod p
   Sy := (twoSy / 2) % p
   fmt.Println(Sy)
}
