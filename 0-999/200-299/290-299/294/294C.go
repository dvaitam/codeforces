package main

import (
   "fmt"
   "sort"
)

const MOD = 1000000007

// modpow computes a^e mod MOD
func modpow(a, e int64) int64 {
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
   var n, m int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   pos := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Scan(&pos[i])
   }
   sort.Ints(pos)

   totalOff := n - m
   // Precompute factorials and inverse factorials up to n
   fact := make([]int64, n+1)
   invFact := make([]int64, n+1)
   fact[0] = 1
   for i := 1; i <= n; i++ {
       fact[i] = fact[i-1] * int64(i) % MOD
   }
   invFact[n] = modpow(fact[n], MOD-2)
   for i := n; i >= 1; i-- {
       invFact[i-1] = invFact[i] * int64(i) % MOD
   }

   ans := fact[totalOff]
   // first block (before the first on-light)
   first := pos[0] - 1
   ans = ans * invFact[first] % MOD
   // internal blocks (between on-lights)
   for i := 1; i < m; i++ {
       gap := pos[i] - pos[i-1] - 1
       ans = ans * invFact[gap] % MOD
       if gap > 0 {
           ans = ans * modpow(2, int64(gap-1)) % MOD
       }
   }
   // last block (after the last on-light)
   last := n - pos[m-1]
   ans = ans * invFact[last] % MOD

   fmt.Println(ans)
