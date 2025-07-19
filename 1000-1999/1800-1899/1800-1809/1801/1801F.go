package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

type pair struct { p, to int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   fmt.Fscan(reader, &n, &k)
   k--
   // Build state intervals for k
   st := make([]int, 0)
   st2 := make([]int, 0)
   id := make(map[int]int)
   cnt := 0
   for i := 1; i <= k; {
       j := k / (k / i)
       st = append(st, j)
       cnt++
       id[j] = cnt
       st2 = append(st2, i)
       i = j + 1
   }
   // dp states: 1..cnt, dp[0] is result
   dp := make([]float64, cnt+1)
   dp[cnt] = float64(k + 1)
   // Transitions for each state
   goTrans := make([][]pair, cnt+1)
   for u := 1; u <= cnt; u++ {
       v := st[u-1]
       capc := int((math.Sqrt(float64(v)) + 1) * 2)
       tmp := make([]pair, 0, capc)
       for i := 2; i <= v; {
           j := v / (v / i)
           tmp = append(tmp, pair{i, id[v/i]})
           i = j + 1
       }
       goTrans[u] = tmp
   }
   // Process each x
   for t := 0; t < n; t++ {
       var x int
       fmt.Fscan(reader, &x)
       fx := float64(x)
       rem := make(map[int]float64, len(st2))
       for i := 1; i <= cnt; i++ {
           d := st2[i-1]
           rem[d] = float64(x/d) / fx
       }
       for u := 1; u <= cnt; u++ {
           pu := dp[u]
           if pu < 1e-8 {
               continue
           }
           // transitions
           for _, pr := range goTrans[u] {
               p, to := pr.p, pr.to
               if p > x {
                   break
               }
               cand := pu * rem[p]
               if cand > dp[to] {
                   dp[to] = cand
               }
           }
           // extra transition when x > v
           v := st[u-1]
           if x > v {
               cand := pu * float64(x/(v+1)) / fx
               if cand > dp[0] {
                   dp[0] = cand
               }
           }
       }
   }
   fmt.Fprintf(writer, "%.10f\n", dp[0])
}
