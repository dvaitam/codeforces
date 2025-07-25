package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)
   // dp[i][l]: number of distinct subsequences of length l using first i chars
   dp := make([][]*big.Int, n+1)
   for i := 0; i <= n; i++ {
       dp[i] = make([]*big.Int, n+1)
       for l := 0; l <= n; l++ {
           dp[i][l] = new(big.Int)
       }
   }
   // empty subsequence
   dp[0][0].SetInt64(1)
   prev := make(map[byte]int)
   for i := 1; i <= n; i++ {
       c := s[i-1]
       // copy counts from i-1
       for l := 0; l <= n; l++ {
           dp[i][l].Set(dp[i-1][l])
       }
       // update by taking c
       for l := 1; l <= i; l++ {
           // add dp[i-1][l-1]
           dp[i][l].Add(dp[i][l], dp[i-1][l-1])
           // subtract duplicates if c appeared before
           if p, ok := prev[c]; ok {
               dp[i][l].Sub(dp[i][l], dp[p-1][l-1])
           }
       }
       // empty subseq always 1
       dp[i][0].SetInt64(1)
       prev[c] = i
   }
   // greedy pick
   rem := k
   sumLen := int64(0)
   for l := n; l >= 0 && rem > 0; l-- {
       cnt := dp[n][l]
       // if cnt >= rem
       take := 0
       if cnt.Cmp(big.NewInt(int64(rem))) >= 0 {
           take = rem
       } else {
           take = int(cnt.Int64())
       }
       sumLen += int64(take) * int64(l)
       rem -= take
   }
   if rem > 0 {
       fmt.Println(-1)
   } else {
       // total cost = k*n - sumLen
       cost := int64(k)*int64(n) - sumLen
       fmt.Println(cost)
   }
}
