package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 998244353

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   freq := make([]int, n+1)
   for _, v := range a {
       freq[v]++
   }
   // precompute inverses
   inv := make([]int, n+1)
   if n >= 1 {
       inv[1] = 1
   }
   for i := 2; i <= n; i++ {
       inv[i] = MOD - int(int64(MOD/i)*int64(inv[MOD%i])%MOD)
   }
   // invFall[k] = 1 / (n * (n-1) * ... * (n-k+1))
   invFall := make([]int, n+1)
   invFall[0] = 1
   for k := 1; k <= n; k++ {
       invFall[k] = int(int64(invFall[k-1]) * int64(inv[n-k+1]) % MOD)
   }
   // dp[s]: sum of products of freq over subsets of size s from processed values
   dp := make([]int, n+1)
   dp[0] = 1
   sz := 0
   ans := 0
   // process values in increasing order
   for v := 1; v <= n; v++ {
       c := freq[v]
       if c >= 2 {
           inner := 0
           // sum over subset sizes s
           for s := 0; s <= sz; s++ {
               if s+2 <= n {
                   inner = (inner + int(int64(dp[s]) * int64(invFall[s+2]) % MOD)) % MOD
               }
           }
           ans = (ans + int(int64(c) * int64(c-1) % MOD * int64(inner) % MOD)) % MOD
       }
       if c > 0 {
           // update dp for including value v
           for s := sz + 1; s >= 1; s-- {
               dp[s] = (dp[s] + int(int64(dp[s-1])*int64(c)%MOD)) % MOD
           }
           sz++
       }
   }
   fmt.Fprintln(writer, ans)
}
