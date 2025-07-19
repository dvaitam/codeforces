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

func abs(a, b int) int {
   if a < b {
       return b - a
   }
   return a - b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &a[i][j])
       }
   }
   // Precompute minimum pairwise differences per column
   minimum := make([][]int, n)
   for i := 0; i < n; i++ {
       minimum[i] = make([]int, n)
   }
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           tmp := int(1e18)
           for k := 0; k < m; k++ {
               d := abs(a[i][k], a[j][k])
               if d < tmp {
                   tmp = d
               }
           }
           minimum[i][j] = tmp
           minimum[j][i] = tmp
       }
   }
   fullMask := (1 << n) - 1
   // dp[fi][pre][mask] = best value
   dp := make([][][]int, n)
   for i := 0; i < n; i++ {
       dp[i] = make([][]int, n)
       for j := 0; j < n; j++ {
           dp[i][j] = make([]int, 1<<n)
           for k := range dp[i][j] {
               dp[i][j][k] = -1
           }
       }
   }
   var calc func(fi, pre, mask int) int
   calc = func(fi, pre, mask int) int {
       if mask == fullMask {
           // compute last edge from pre to fi
           r := int(1e18)
           for i := 0; i < m-1; i++ {
               d := abs(a[pre][i], a[fi][i+1])
               if d < r {
                   r = d
               }
           }
           return r
       }
       if dp[fi][pre][mask] != -1 {
           return dp[fi][pre][mask]
       }
       bt := 0
       for i := 0; i < n; i++ {
           if mask&(1<<i) == 0 {
               c := calc(fi, i, mask|(1<<i))
               if minimum[pre][i] < c {
                   c = minimum[pre][i]
               }
               if c > bt {
                   bt = c
               }
           }
       }
       dp[fi][pre][mask] = bt
       return bt
   }
   res := 0
   for i := 0; i < n; i++ {
       v := calc(i, i, 1<<i)
       if v > res {
           res = v
       }
   }
   fmt.Fprintln(writer, res)
}
