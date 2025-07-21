package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type marble struct {
   x, c int64
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   arr := make([]marble, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &arr[i].x, &arr[i].c)
   }
   sort.Slice(arr, func(i, j int) bool { return arr[i].x < arr[j].x })
   x := make([]int64, n+1)
   c := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       x[i] = arr[i-1].x
       c[i] = arr[i-1].c
   }
   prefixX := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       prefixX[i] = prefixX[i-1] + x[i]
   }
   const INF = int64(9e18)
   dp := make([]int64, n+1)
   // Must pin the first marble to avoid infinite roll
   dp[1] = c[1]
   // DP: dp[i] = min cost to cover up to i and pin i
   for i := 2; i <= n; i++ {
       best := INF
       for p := 1; p <= i-1; p++ {
           // marbles p+1..i-1 roll to x[p]
           moveCost := (prefixX[i-1] - prefixX[p]) - int64(i-1-p)*x[p]
           if dp[p]+moveCost < best {
               best = dp[p] + moveCost
           }
       }
       dp[i] = best + c[i]
   }
   // Consider tail: marbles after last pinned roll to it
   ans := INF
   for i := 1; i <= n; i++ {
       tailCost := (prefixX[n] - prefixX[i]) - int64(n-i)*x[i]
       if dp[i]+tailCost < ans {
           ans = dp[i] + tailCost
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
}
