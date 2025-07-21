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
   if k > n {
       fmt.Println(0)
       return
   }
   // Precompute powers of 2 up to n*n
   maxExp := n * n
   pow2 := make([]int, maxExp+1)
   pow2[0] = 1
   for i := 1; i <= maxExp; i++ {
       pow2[i] = pow2[i-1] * 2 % mod
   }
   // Precompute transition A[prev][i] for 0<=prev<i<=n
   A := make([][]int, n+1)
   for prev := 0; prev <= n; prev++ {
       A[prev] = make([]int, n+1)
       for i := prev + 1; i <= n; i++ {
           d := i - prev
           // (2^d - 1) * 2^(d*prev)
           v := pow2[d] - 1
           if v < 0 {
               v += mod
           }
           // exponent = d * prev
           e := d * prev
           v = int(int64(v) * int64(pow2[e]) % mod)
           A[prev][i] = v
       }
   }
   // dp[t][i]: sequences length t ending at i
   dp := make([][]int, k+1)
   for t := 0; t <= k; t++ {
       dp[t] = make([]int, n+1)
   }
   dp[0][0] = 1
   for t := 1; t <= k; t++ {
       // for each possible current end i
       for i := 1; i <= n; i++ {
           var sum int64
           // prev from t-1 picks must be at least t-1, so prev >= t-1
           // but dp[t-1][prev] is zero for prev < t-1 (except prev=0 when t=1)
           // We'll just loop prev from 0..i-1
           for prev := 0; prev < i; prev++ {
               if dp[t-1][prev] != 0 {
                   sum += int64(dp[t-1][prev]) * int64(A[prev][i])
                   if sum >= 1<<63 {
                       sum %= mod
                   }
               }
           }
           dp[t][i] = int(sum % mod)
       }
   }
   // Sum up dp[k][i] * 2^((n-i)*i)
   var ans int64
   for i := 0; i <= n; i++ {
       if dp[k][i] == 0 {
           continue
       }
       // exponent for final gap
       e := (n - i) * i
       ans = (ans + int64(dp[k][i])*int64(pow2[e])) % mod
   }
   fmt.Println(ans)
}
