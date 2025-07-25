package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 998244353

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, c int
   fmt.Fscan(reader, &n, &c)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // max full rounds possible is floor(n/c), add extra for safety
   maxR := n/c + 2
   // dp[j][t]: number of subsequences with j full rounds completed, waiting for symbol t (1..c)
   dp := make([][]int, maxR)
   for j := range dp {
       dp[j] = make([]int, c+2)
   }
   dp[0][1] = 1
   // process each element
   for _, x := range a {
       // new dp
       newdp := make([][]int, maxR)
       for j := range newdp {
           newdp[j] = make([]int, c+2)
       }
       for j := 0; j < maxR; j++ {
           for t := 1; t <= c; t++ {
               v := dp[j][t]
               if v == 0 {
                   continue
               }
               // skip current element
               nd := (newdp[j][t] + v) % MOD
               newdp[j][t] = nd
               // include current element
               if x == t {
                   if t < c {
                       newdp[j][t+1] = (newdp[j][t+1] + v) % MOD
                   } else {
                       // complete a round
                       if j+1 < maxR {
                           newdp[j+1][1] = (newdp[j+1][1] + v) % MOD
                       }
                   }
               } else {
                   // include does not change state
                   newdp[j][t] = (newdp[j][t] + v) % MOD
               }
           }
       }
       dp = newdp
   }
   // accumulate answers
   // s_p: number of subsequences with exactly p full rounds
   res := make([]int, n+1)
   for j := 0; j < maxR; j++ {
       for t := 1; t <= c; t++ {
           v := dp[j][t]
           if v != 0 {
               if j <= n {
                   res[j] = (res[j] + v) % MOD
               }
           }
       }
   }
   // exclude empty subsequence from res[0]
   res[0] = (res[0] - 1 + MOD) % MOD
   // print s_0 ... s_n
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for p := 0; p <= n; p++ {
       if p > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprint(res[p]))
   }
   writer.WriteByte('\n')
}
