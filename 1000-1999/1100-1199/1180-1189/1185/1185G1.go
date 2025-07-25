package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

const mod = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, T int
   fmt.Fscan(reader, &n, &T)
   t := make([]int, n)
   g := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &t[i], &g[i])
   }
   size := 1 << n
   // Precompute total time for each subset
   timeSum := make([]int, size)
   for mask := 1; mask < size; mask++ {
       lb := mask & -mask
       idx := bits.TrailingZeros(uint(lb))
       timeSum[mask] = timeSum[mask^lb] + t[idx]
   }
   // dp[mask][lastGenre] = number of ways
   dp := make([][4]int, size)
   dp[0][0] = 1
   for mask := 0; mask < size; mask++ {
       for last := 0; last < 4; last++ {
           val := dp[mask][last]
           if val == 0 {
               continue
           }
           for i := 0; i < n; i++ {
               if mask&(1<<i) != 0 || g[i] == last {
                   continue
               }
               nm := mask | (1 << i)
               dp[nm][g[i]] += val
               if dp[nm][g[i]] >= mod {
                   dp[nm][g[i]] -= mod
               }
           }
       }
   }
   ans := 0
   for mask := 0; mask < size; mask++ {
       if timeSum[mask] != T {
           continue
       }
       for last := 1; last <= 3; last++ {
           ans += dp[mask][last]
           if ans >= mod {
               ans -= mod
           }
       }
   }
   fmt.Fprint(writer, ans)
}
