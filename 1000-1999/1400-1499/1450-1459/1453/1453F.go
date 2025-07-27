package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
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
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n+1)
       // a[1..n]
       for i := 1; i <= n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       // dp[u][v]: max nodes from start (1) to v inclusive, where previous node is u
       // we index u from 0..n, v from 1..n
       inf := -1_000_000_000
       dp := make([][]int, n+1)
       for i := 0; i <= n; i++ {
           dp[i] = make([]int, n+1)
           for j := 0; j <= n; j++ {
               dp[i][j] = inf
           }
       }
       // base: before start, use u=0 with a[0]=0
       dp[0][1] = 1
       // process
       B := make([]int, n+1)
       pref := make([]int, n+1)
       for v := 1; v < n; v++ {
           // prepare B
           for i := 0; i <= n; i++ {
               B[i] = inf
           }
           for u := 0; u < v; u++ {
               val := dp[u][v]
               if val > inf {
                   h := u
                   if u <= n {
                       h = u + aVal(a, u)
                   }
                   if h > n {
                       h = n
                   }
                   if B[h] < val {
                       B[h] = val
                   }
               }
           }
           // prefix max: pref[k] = max B[0..k-1]
           pref[0] = inf
           for k := 1; k <= n; k++ {
               pref[k] = pref[k-1]
               if B[k-1] > pref[k] {
                   pref[k] = B[k-1]
               }
           }
           // transitions v->w
           for w := v + 1; w <= n && w <= v+a[v]; w++ {
               // dp[v][w] = max(dp[v][w], pref[w] + 1)
               if pref[w] > inf {
                   dp[v][w] = max(dp[v][w], pref[w]+1)
               }
           }
       }
       // find best ending at n
       best := 0
       for u := 0; u < n; u++ {
           if dp[u][n] > best {
               best = dp[u][n]
           }
       }
       // answer: zero count = n - best
       fmt.Fprintln(writer, n-best)
   }
}

// aVal returns a[u] for u>=1, or 0 for u==0
func aVal(a []int, u int) int {
   if u >= 1 && u < len(a) {
       return a[u]
   }
   return 0
}
