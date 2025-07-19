package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   d := make([]int, 2*n+2)
   c := make([]int, 2*n+2)
   for i := 1; i <= 2*n; i++ {
       fmt.Fscan(in, &d[i])
       c[d[i]] = i
   }
   b := make([]int, 2*n+2)
   type pair struct{ val, length int }
   var v []pair
   idx := 2 * n
   for i := 2 * n; i >= 1; i-- {
       if c[i] > idx {
           continue
       }
       length := idx - c[i] + 1
       v = append(v, pair{i, length})
       b[i] = length
       idx = c[i] - 1
   }
   // reverse v
   for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
       v[i], v[j] = v[j], v[i]
   }
   m := len(v)
   // dp[i][j]: using first i+1 items, sum j possible
   dp := make([][]bool, m)
   for i := 0; i < m; i++ {
       dp[i] = make([]bool, n+1)
       dp[i][0] = true
   }
   if m > 0 && v[0].length <= n {
       dp[0][v[0].length] = true
   }
   idxRes := -1
   if m > 0 && dp[0][n] {
       idxRes = 0
   }
   for i := 1; i < m; i++ {
       for j := 1; j <= n; j++ {
           if dp[i-1][j] {
               dp[i][j] = true
           } else if j >= v[i].length && dp[i-1][j-v[i].length] {
               dp[i][j] = true
           }
       }
       if dp[i][n] {
           idxRes = i
       }
   }
   if idxRes == -1 {
       fmt.Fprintln(out, -1)
       return
   }
   // backtrack
   x := idxRes
   y := n
   ansSum := 0
   var vv []int
   for ansSum != n {
       if x == 0 {
           vv = append(vv, v[0].val)
           break
       }
       if dp[x-1][y] {
           x--
       } else {
           vv = append(vv, v[x].val)
           ansSum += v[x].length
           y -= v[x].length
           x--
       }
   }
   // reverse vv
   for i, j := 0, len(vv)-1; i < j; i, j = i+1, j-1 {
       vv[i], vv[j] = vv[j], vv[i]
   }
   used := make([]bool, 2*n+2)
   for _, val := range vv {
       pos := c[val]
       length := b[val]
       for j := pos; j < pos+length; j++ {
           fmt.Fprint(out, d[j], " ")
           used[j] = true
       }
   }
   fmt.Fprintln(out)
   for i := 1; i <= 2*n; i++ {
       if !used[i] {
           fmt.Fprint(out, d[i], " ")
       }
   }
   fmt.Fprintln(out)
}
