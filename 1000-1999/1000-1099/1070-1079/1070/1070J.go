package main

import (
   "bufio"
   "fmt"
   "os"
)

func minInt64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
       var n, m int64
       var k int
       fmt.Fscan(reader, &n, &m, &k)
       var s string
       fmt.Fscan(reader, &s)

       const ALPH = 26
       cnt := make([]int, ALPH)
       for _, ch := range s {
           cnt[ch-'A']++
       }

       dp := make([]int64, k+1)
       dp[0] = 1
       // build subset count
       for _, x := range cnt {
           if x == 0 {
               continue
           }
           for j := k; j >= x; j-- {
               dp[j] += dp[j-x]
           }
       }

       ans := int64(k) * int64(k)
       // iterate each character count
       for _, x := range cnt {
           if x == 0 {
               continue
           }
           // remove current letter contributions
           for j := x; j <= k; j++ {
               dp[j] -= dp[j-x]
           }

           for i := 0; i <= k; i++ {
               if dp[i] == 0 {
                   continue
               }
               rm := k - i - x
               n1 := n - minInt64(n, int64(i))
               m1 := m - minInt64(m, int64(rm))
               if n1 + m1 > int64(x) {
                   continue
               }
               prod := n1 * m1
               if prod < ans {
                   ans = prod
               }
           }
           // restore dp
           for j := k; j >= x; j-- {
               dp[j] += dp[j-x]
           }
       }

       fmt.Fprintln(writer, ans)
   }
}
