package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, s, k int
   if _, err := fmt.Fscan(reader, &n, &s, &k); err != nil {
       return
   }
   s-- // zero-index
   r := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &r[i])
   }
   var c string
   fmt.Fscan(reader, &c)

   const INF = 1 << 60
   best := INF
   // dp[i][sum] = min cost to reach i with accumulated rating sum
   dp := make([][]int, n)
   for i := 0; i < n; i++ {
       dp[i] = make([]int, k)
       for j := range dp[i] {
           dp[i][j] = INF
       }
   }
   // initial positions
   for i := 0; i < n; i++ {
       cost := abs(i - s)
       if r[i] >= k {
           best = min(best, cost)
       } else {
           dp[i][r[i]] = cost
       }
   }
   // transitions by increasing rating
   for R := 1; R <= 50; R++ {
       for i := 0; i < n; i++ {
           if r[i] != R {
               continue
           }
           for sum := 0; sum < k; sum++ {
               curr := dp[i][sum]
               if curr >= INF {
                   continue
               }
               for j := 0; j < n; j++ {
                   if r[j] > r[i] && c[i] != c[j] {
                       s2 := sum + r[j]
                       cost2 := curr + abs(j-i)
                       if s2 >= k {
                           best = min(best, cost2)
                       } else {
                           if cost2 < dp[j][s2] {
                               dp[j][s2] = cost2
                           }
                       }
                   }
               }
           }
       }
   }
   if best < INF {
       fmt.Println(best)
   } else {
       fmt.Println(-1)
   }
}
