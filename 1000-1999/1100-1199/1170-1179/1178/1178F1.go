package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 998244353

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   c := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &c[i])
   }
   // Precompute min positions for all ranges
   minPos := make([][]int, n+2)
   for i := 1; i <= n; i++ {
       minPos[i] = make([]int, n+1)
       minPos[i][i] = i
       for j := i + 1; j <= n; j++ {
           prev := minPos[i][j-1]
           if c[j] < c[prev] {
               minPos[i][j] = j
           } else {
               minPos[i][j] = prev
           }
       }
   }
   // dp[i][j]: number of ways for segment [i..j]
   dp := make([][]int, n+3)
   for i := 1; i <= n+2; i++ {
       dp[i] = make([]int, n+2)
   }
   // base: dp[i][i-1] = 1
   for i := 1; i <= n+1; i++ {
       dp[i][i-1] = 1
   }
   // compute dp by increasing length
   for length := 1; length <= n; length++ {
       for i := 1; i+length-1 <= n; i++ {
           j := i + length - 1
           p := minPos[i][j]
           // compute left sum
           var leftSum int
           for a := i; a <= p; a++ {
               left := dp[i][a-1]
               right := dp[a][p-1]
               leftSum = (leftSum + left*right) % MOD
           }
           // compute right sum
           var rightSum int
           for b := p; b <= j; b++ {
               left := dp[p+1][b]
               right := dp[b+1][j]
               rightSum = (rightSum + left*right) % MOD
           }
           dp[i][j] = leftSum * rightSum % MOD
       }
   }
   // output answer for [1..n]
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, dp[1][n])
}
