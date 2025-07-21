package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func add(a, b int) int {
   a += b
   if a >= MOD {
       a -= MOD
   }
   return a
}

func mul(a, b int) int {
   return int((int64(a) * int64(b)) % MOD)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   fmt.Fscan(reader, &n, &k)
   var s string
   fmt.Fscan(reader, &s)

   // dp[l][b]: number of ways at current position with run length l, beauty b
   dp := make([][]int, n+1)
   for i := range dp {
       dp[i] = make([]int, k+1)
   }
   dp[0][0] = 1

   for r := 1; r <= n; r++ {
       // position r, character s[r-1]
       ch := s[r-1]
       cntLt := int(ch - 'a')
       cntGt := 26 - cntLt - 1
       m := n - r + 1
       // prefix sums for < transitions
       sumB := make([]int, k+1)
       for l := 0; l < r; l++ {
           for b := 0; b <= k; b++ {
               if dp[l][b] != 0 {
                   sumB[b] = add(sumB[b], dp[l][b])
               }
           }
       }
       // new dp
       ndp := make([][]int, n+1)
       for i := range ndp {
           ndp[i] = make([]int, k+1)
       }
       // transitions
       for b := 0; b <= k; b++ {
           if sumB[b] != 0 && cntLt > 0 {
               ndp[0][b] = add(ndp[0][b], mul(sumB[b], cntLt))
           }
       }
       // eq and > transitions
       for l := 0; l < r; l++ {
           for b := 0; b <= k; b++ {
               v := dp[l][b]
               if v == 0 {
                   continue
               }
               // eq
               ndp[l+1][b] = add(ndp[l+1][b], v)
               // >
               if cntGt > 0 {
                   delta := (l + 1) * m
                   nb := b + delta
                   if nb <= k {
                       ndp[0][nb] = add(ndp[0][nb], mul(v, cntGt))
                   }
               }
           }
       }
       dp = ndp
   }
   // sum all run lengths for beauty k
   ans := 0
   for l := 0; l <= n; l++ {
       ans = add(ans, dp[l][k])
   }
   fmt.Fprintln(writer, ans)
}
