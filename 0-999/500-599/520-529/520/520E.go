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
   // nCk helper
   comb := func(n, k int) int64 {
       if k < 0 || k > n || n < 0 {
           return 0
       }
       return fact[n] * invf[k] % MOD * invf[n-k] % MOD
   }
   // precompute powers of 10
   pow10 := make([]int64, n+2)
   pow10[0] = 1
   for i := 1; i <= n+1; i++ {
       pow10[i] = pow10[i-1] * 10 % MOD
   }
   // build prefix sums of f[d] * 10^d, where f[d] = C(n-2-d, k-1)
   maxd := n - 1
   S := make([]int64, maxd)
   for d := 0; d < maxd; d++ {
       val := comb(n-2-d, k-1) * pow10[d] % MOD
       if d == 0 {
           S[d] = val
       } else {
           S[d] = (S[d-1] + val) % MOD
       }
   }
   var ans int64
   for i := 1; i <= n; i++ {
       d := int64(s[i-1] - '0')
       if d == 0 {
           continue
       }
       // sum over segments ending before n: d from 0 to n-i-1
       var sum1 int64
       rem := n - i - 1
       if rem >= 0 {
           sum1 = S[rem]
       }
       // segments ending at n: one segment tail
       sum2 := comb(i-1, k) * pow10[n-i] % MOD
       total := (sum1 + sum2) % MOD
       ans = (ans + d*total) % MOD
   }
   fmt.Println(ans)
}
