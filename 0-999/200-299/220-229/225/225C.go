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

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, x, y int
   if _, err := fmt.Fscan(reader, &n, &m, &x, &y); err != nil {
       return
   }
   // count black pixels in each column
   blacks := make([]int, m+1)
   for i := 0; i < n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       for j := 0; j < m; j++ {
           if line[j] == '#' {
               blacks[j+1]++
           }
       }
   }
   // cost to paint each column white or black
   costWhite := make([]int, m+1)
   costBlack := make([]int, m+1)
   for j := 1; j <= m; j++ {
       costWhite[j] = blacks[j]
       costBlack[j] = n - blacks[j]
   }
   // prefix sums
   prefWhite := make([]int, m+1)
   prefBlack := make([]int, m+1)
   for j := 1; j <= m; j++ {
       prefWhite[j] = prefWhite[j-1] + costWhite[j]
       prefBlack[j] = prefBlack[j-1] + costBlack[j]
   }
   const INF = 1e9
   // dp[j][c]: min cost for first j columns ending with color c (0 white,1 black)
   dp := make([][2]int, m+1)
   for j := 0; j <= m; j++ {
       dp[j][0], dp[j][1] = INF, INF
   }
   dp[0][0], dp[0][1] = 0, 0
   for j := 1; j <= m; j++ {
       // ending white
       for k := x; k <= y && k <= j; k++ {
           cost := prefWhite[j] - prefWhite[j-k]
           prev := dp[j-k][1]
           if prev+cost < dp[j][0] {
               dp[j][0] = prev + cost
           }
       }
       // ending black
       for k := x; k <= y && k <= j; k++ {
           cost := prefBlack[j] - prefBlack[j-k]
           prev := dp[j-k][0]
           if prev+cost < dp[j][1] {
               dp[j][1] = prev + cost
           }
       }
   }
   // answer
   ans := dp[m][0]
   if dp[m][1] < ans {
       ans = dp[m][1]
   }
   fmt.Println(ans)
}
