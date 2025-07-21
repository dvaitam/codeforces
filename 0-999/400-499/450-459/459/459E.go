package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type edge struct {
   u, v, w int
}

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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   edges := make([]edge, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &edges[i].u, &edges[i].v, &edges[i].w)
   }
   sort.Slice(edges, func(i, j int) bool {
       return edges[i].w < edges[j].w
   })

   dp := make([]int, n+1)
   temp := make([]int, n+1)
   vs := make([]int, 0)
   ans := 0
   for i := 0; i < m; {
       j := i
       // group edges with same weight
       for j < m && edges[j].w == edges[i].w {
           j++
       }
       vs = vs[:0]
       // compute temp dp for this weight
       for k := i; k < j; k++ {
           u := edges[k].u
           v := edges[k].v
           cand := dp[u] + 1
           if cand > temp[v] {
               if temp[v] == 0 {
                   vs = append(vs, v)
               }
               temp[v] = cand
           }
       }
       // apply updates
       for _, v := range vs {
           dp[v] = max(dp[v], temp[v])
           ans = max(ans, dp[v])
           temp[v] = 0
       }
       i = j
   }

   fmt.Fprintln(writer, ans)
}
