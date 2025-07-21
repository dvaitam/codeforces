package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func modPow(a, e int64) int64 {
   res := int64(1)
   for e > 0 {
       if e&1 == 1 {
           res = res * a % MOD
       }
       a = a * a % MOD
       e >>= 1
   }
   return res
}

func modInv(a int64) int64 {
   // MOD is prime
   return modPow(a, MOD-2)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   s := make([]byte, n)
   fmt.Fscan(in, &s)
   // special case k==0: whole string
   if k == 0 {
       var ans int64
       for i := 0; i < n; i++ {
           ans = (ans*10 + int64(s[i]-'0')) % MOD
       }
       fmt.Println(ans)
       return
   }
   // precompute factorials and invfacts
   N := n
   fact := make([]int64, N+1)
   invf := make([]int64, N+1)
   fact[0] = 1
   for i := 1; i <= N; i++ {
       fact[i] = fact[i-1] * int64(i) % MOD
   }
   invf[N] = modInv(fact[N])
   for i := N; i > 0; i-- {
       invf[i-1] = invf[i] * int64(i) % MOD
   }
   // nCk function
   comb := func(n, k int) int64 {
       if k < 0 || k > n || n < 0 {
           return 0
       }
       return fact[n] * invf[k] % MOD * invf[n-k] % MOD
   }
   // precompute powers of 10
   pow10 := make([]int64, n+1)
   pow10[0] = 1
   for i := 1; i <= n; i++ {
       pow10[i] = pow10[i-1] * 10 % MOD
   }
   inv9 := modInv(9)
   // precompute reusable combs
   c_n2 := comb(n-2, k-1)
   c_n3 := comb(n-3, k-2)
   var ans int64
   // iterate positions 1..n as i index 0..n-1
   for i := 0; i < n; i++ {
       d := int64(s[i] - '0')
       if d == 0 {
           continue
       }
       // A = 10^(n-1-i)
       exp := int64(n-1-i)
       A := pow10[n-1-i]
       // B = (10^(n-i) - 1) / 9 = sum_{d=0..n-1-i} 10^d
       B := (pow10[n-i] + MOD - 1) % MOD * inv9 % MOD
       // term_start: l=1, r from i..n-1: C(n-2,k-1) * B
       term := d * c_n2 % MOD * B % MOD
       ans = (ans + term) % MOD
       // term_mid: l>1, r<n: (i) choices of l (positions 2..i+1 => count=i) but index i gives i choices
       if c_n3 != 0 {
           term = d * c_n3 % MOD * int64(i) % MOD * B % MOD
           ans = (ans + term) % MOD
       }
       // term_end: r=n, l>1: C(n-2,k-1) * (i) * A
       term = d * c_n2 % MOD * int64(i) % MOD * A % MOD
       ans = (ans + term) % MOD
   }
   fmt.Println(ans)
}
