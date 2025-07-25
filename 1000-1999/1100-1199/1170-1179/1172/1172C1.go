package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353

func modPow(a, e int) int {
   res := 1
   for e > 0 {
       if e&1 != 0 {
           res = res * a % mod
       }
       a = a * a % mod
       e >>= 1
   }
   return res
}

func modInv(a int) int {
   // mod prime
   return modPow(a, mod-2)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(in, &n, &m)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   w := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &w[i])
   }
   sum0, sum1 := 0, 0
   for i := 0; i < n; i++ {
       if a[i] == 1 {
           sum1 += w[i]
       } else {
           sum0 += w[i]
       }
   }
   // dp[x][y]: probability to have x liked, y disliked picks after x+y steps
   // only need for x+y <= m
   dp := make([][]int, m+1)
   for i := range dp {
       dp[i] = make([]int, m+1)
   }
   dp[0][0] = 1
   tot := sum0 + sum1
   // precompute inverses
   maxDen := tot + m
   inv := make([]int, maxDen+1)
   for i := 1; i <= maxDen; i++ {
       inv[i] = modInv(i)
   }
   for x := 0; x <= m; x++ {
       for y := 0; x+y < m; y++ {
           cur := dp[x][y]
           if cur == 0 {
               continue
           }
           rem := tot + x - y
           invRem := inv[rem]
           // liked
           num1 := sum1 + x
           add1 := cur * num1 % mod * invRem % mod
           dp[x+1][y] = (dp[x+1][y] + add1) % mod
           // disliked
           if sum0-y > 0 {
               num0 := sum0 - y
               add0 := cur * num0 % mod * invRem % mod
               dp[x][y+1] = (dp[x][y+1] + add0) % mod
           }
       }
   }
   // expected total liked picks E[X]
   expX := 0
   for x := 0; x <= m; x++ {
       y := m - x
       if y < 0 || y > m {
           continue
       }
       expX = (expX + dp[x][y]*x) % mod
   }
   // E[Y] = m - E[X]
   expY := (m%mod - expX + mod) % mod
   // prepare answers
   invS1 := 0
   if sum1 > 0 {
       invS1 = modInv(sum1)
   }
   invS0 := 0
   if sum0 > 0 {
       invS0 = modInv(sum0)
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i := 0; i < n; i++ {
       wi := w[i]
       var ans int
       if a[i] == 1 {
           // wi * (sum1 + expX) / sum1
           tmp := (sum1 + expX) % mod
           ans = wi % mod * tmp % mod * invS1 % mod
       } else {
           // wi * (sum0 - expY) / sum0
           tmp := (sum0 - expY) % mod
           if tmp < 0 {
               tmp += mod
           }
           ans = wi % mod * tmp % mod * invS0 % mod
       }
       fmt.Fprint(out, ans)
       if i+1 < n {
           fmt.Fprint(out, " ")
       }
   }
   fmt.Fprintln(out)
}
