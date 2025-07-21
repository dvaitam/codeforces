package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   const INF = int64(1<<62)
   // cost2[l][r]: minimal cost to fix segment [l,r]
   cost2 := make([][]int64, n+2)
   for l := 1; l <= n; l++ {
       cost2[l] = make([]int64, n+2)
       for r := 1; r <= n; r++ {
           cost2[l][r] = INF
       }
   }
   for i := 0; i < m; i++ {
       var l, r int
       var c int64
       fmt.Fscan(reader, &l, &r, &c)
       if c < cost2[l][r] {
           cost2[l][r] = c
       }
   }
   // minCostFrom[a][x]: minimal cost among segments starting at a and ending at >= x
   minCostFrom := make([][]int64, n+2)
   for a := 1; a <= n; a++ {
       minCostFrom[a] = make([]int64, n+3)
       // copy direct costs
       for x := 1; x <= n; x++ {
           minCostFrom[a][x] = cost2[a][x]
       }
       minCostFrom[a][n+1] = INF
       // propagate minima backwards
       for x := n; x >= 1; x-- {
           if minCostFrom[a][x] > minCostFrom[a][x+1] {
               minCostFrom[a][x] = minCostFrom[a][x+1]
           }
       }
   }
   ans := INF
   // f[x]: minimal cost to cover all holes in [L..x]
   f := make([]int64, n+2)
   for L := 1; L <= n; L++ {
       // init f
       for i := L - 1; i <= n; i++ {
           f[i] = INF
       }
       f[L-1] = 0
       // DP to cover [L..x]
       for x := L; x <= n; x++ {
           var best int64 = INF
           for a := L; a <= x; a++ {
               if f[a-1] == INF {
                   continue
               }
               c := minCostFrom[a][x]
               if c == INF {
                   continue
               }
               v := f[a-1] + c
               if v < best {
                   best = v
               }
           }
           f[x] = best
       }
       // any block of length >= k starting at L
       end := L + k - 1
       if end > n {
           continue
       }
       for x := end; x <= n; x++ {
           if f[x] < ans {
               ans = f[x]
           }
       }
   }
   if ans >= INF {
       fmt.Fprintln(writer, -1)
   } else {
       fmt.Fprintln(writer, ans)
   }
}
