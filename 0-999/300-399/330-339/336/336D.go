package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

// modpow computes a^e mod mod
func modpow(a, e int64) int64 {
   var res int64 = 1
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

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   var g int
   fmt.Fscan(reader, &n, &m, &g)
   // Special case: no ones
   if m == 0 {
       var B int
       if n%2 == 0 {
           B = 1
       } else {
           B = 0
       }
       var ans int
       if g == 1 {
           ans = B
       } else {
           ans = (1 - B + mod) % mod
       }
       fmt.Println(ans)
       return
   }
   // Precompute factorials and inverse factorials up to n + m
   N := n + m
   fact := make([]int64, N+1)
   invfact := make([]int64, N+1)
   fact[0] = 1
   for i := 1; i <= N; i++ {
       fact[i] = fact[i-1] * int64(i) % mod
   }
   invfact[N] = modpow(fact[N], mod-2)
   for i := N; i >= 1; i-- {
       invfact[i-1] = invfact[i] * int64(i) % mod
   }
   // comb returns C(a, b)
   comb := func(a, b int) int64 {
       if b < 0 || b > a {
           return 0
       }
       return fact[a] * invfact[b] % mod * invfact[a-b] % mod
   }
   // B[i] = number of strings with i zeros, m ones that reduce to 1
   B := make([]int64, n+1)
   // Base cases
   if m == 1 {
       B[0] = 1
   } else {
       B[0] = 0
   }
   if n >= 1 {
       if m >= 2 {
           B[1] = 1
       } else {
           B[1] = 0
       }
   }
   // Recurrence: B[i] = C(i+m-2, i-1) + B[i-2]
   for i := 2; i <= n; i++ {
       B[i] = (comb(i+m-2, i-1) + B[i-2]) % mod
   }
   total := comb(n+m, n)
   Bn := B[n]
   var ans int64
   if g == 1 {
       ans = Bn
   } else {
       ans = (total - Bn + mod) % mod
   }
   fmt.Println(ans)
