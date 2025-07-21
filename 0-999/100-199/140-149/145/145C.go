package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   freq := make(map[int]int)
   var U int
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(in, &x)
       if isLucky(x) {
           freq[x]++
       } else {
           U++
       }
   }
   // counts of each distinct lucky number
   counts := make([]int, 0, len(freq))
   for _, c := range freq {
       counts = append(counts, c)
   }
   m := len(counts)
   // dp[t] = number of ways to pick t distinct lucky numbers
   dp := make([]int64, m+1)
   dp[0] = 1
   for _, c := range counts {
       cc := int64(c)
       for t := m; t >= 1; t-- {
           dp[t] = (dp[t] + dp[t-1]*cc) % mod
       }
   }
   // precompute factorials up to n
   fact := make([]int64, n+1)
   invFact := make([]int64, n+1)
   fact[0] = 1
   for i := 1; i <= n; i++ {
       fact[i] = fact[i-1] * int64(i) % mod
   }
   invFact[n] = modInv(fact[n])
   for i := n; i > 0; i-- {
       invFact[i-1] = invFact[i] * int64(i) % mod
   }
   // compute answer
   var ans int64
   maxT := m
   if k < maxT {
       maxT = k
   }
   for t := 0; t <= maxT; t++ {
       rem := k - t
       if rem < 0 || rem > U {
           continue
       }
       waysLucky := dp[t]
       waysUnlucky := nCr(int64(U), rem, fact, invFact)
       ans = (ans + waysLucky*waysUnlucky) % mod
   }
   fmt.Println(ans)
}

// isLucky returns true if x consists only of digits 4 and 7
func isLucky(x int) bool {
   if x <= 0 {
       return false
   }
   for x > 0 {
       d := x % 10
       if d != 4 && d != 7 {
           return false
       }
       x /= 10
   }
   return true
}

// modPow computes a^b % mod
func modPow(a, b int64) int64 {
   res := int64(1)
   a %= mod
   for b > 0 {
       if b&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       b >>= 1
   }
   return res
}

// modInv computes modular inverse of a mod mod
func modInv(a int64) int64 {
   return modPow(a, mod-2)
}

// nCr computes C(n, r) using precomputed factorials
func nCr(n int64, r int, fact, invFact []int64) int64 {
   if r < 0 || n < int64(r) {
       return 0
   }
   return fact[n] * invFact[r] % mod * invFact[n-int64(r)] % mod
}
