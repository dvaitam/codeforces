package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 10007

func modpow(a, e int64) int64 {
   res := int64(1)
   base := a % mod
   for e > 0 {
       if e&1 == 1 {
           res = (res * base) % mod
       }
       base = (base * base) % mod
       e >>= 1
   }
   return res
}

var fact, invfact []int64

// compute nCr % mod with Lucas
func comb(n, k int64) int64 {
   if k < 0 || k > n {
       return 0
   }
   if n == 0 || k == 0 {
       return 1
   }
   // Lucas theorem
   ni := n % mod
   ki := k % mod
   if ki > ni {
       return 0
   }
   return comb(n/mod, k/mod) * (fact[ni] * invfact[ki] % mod * invfact[ni-ki] % mod) % mod
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var s string
   var n int64
   fmt.Fscan(in, &s)
   fmt.Fscan(in, &n)
   m := len(s)
   // DP for minimal insertions and count
   dp := make([][]int, m)
   cnt := make([][]int64, m)
   for i := 0; i < m; i++ {
       dp[i] = make([]int, m)
       cnt[i] = make([]int64, m)
       dp[i][i] = 0
       cnt[i][i] = 1
   }
   // lengths
   for length := 2; length <= m; length++ {
       for i := 0; i+length-1 < m; i++ {
           j := i + length - 1
           if s[i] == s[j] {
               // inside
               if i+1 <= j-1 {
                   dp[i][j] = dp[i+1][j-1]
                   cnt[i][j] = cnt[i+1][j-1]
               } else {
                   // empty or single
                   dp[i][j] = 0
                   cnt[i][j] = 1
               }
           } else {
               // insert either s[i] or s[j]
               left := dp[i+1][j]
               right := dp[i][j-1]
               if left < right {
                   dp[i][j] = left + 1
                   cnt[i][j] = cnt[i+1][j]
               } else if left > right {
                   dp[i][j] = right + 1
                   cnt[i][j] = cnt[i][j-1]
               } else {
                   dp[i][j] = left + 1
                   cnt[i][j] = (cnt[i+1][j] + cnt[i][j-1]) % mod
               }
           }
       }
   }
   D := dp[0][m-1]
   cntMin := cnt[0][m-1]
   surplus := n - int64(D)
   if surplus < 0 {
       fmt.Println(0)
       return
   }
   // init factorials
   fact = make([]int64, mod)
   invfact = make([]int64, mod)
   fact[0] = 1
   for i := 1; i < mod; i++ {
       fact[i] = fact[i-1] * int64(i) % mod
   }
   invfact[mod-1] = modpow(int64(fact[mod-1]), mod-2)
   for i := mod - 1; i > 0; i-- {
       invfact[i-1] = invfact[i] * int64(i) % mod
   }
   // minimal palindrome length
   L0 := int64(m + D)
   // number of pairs to insert
   hasMid := surplus & 1
   k := surplus / 2
   // number of slots for pairs
   slots := L0 + 1
   // ways to distribute k pairs into slots: C(k+slots-1, slots-1)
   ways := comb(k+slots-1, slots-1)
   // choose letters for pairs
   ways = ways * modpow(26, k) % mod
   if hasMid == 1 {
       // one center insertion: if P0 length even, any of 26 letters; if odd, only match P0 center (1 choice)
       if L0%2 == 0 {
           ways = ways * 26 % mod
       }
   }
   ans := cntMin * ways % mod
   fmt.Println(ans)
}
