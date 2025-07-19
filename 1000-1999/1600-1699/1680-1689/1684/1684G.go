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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // build bipartite graph: big nodes to all nodes
   adj := make([][]int, n)
   big := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if a[i] > m/3 {
           big = append(big, i)
           for j := 0; j < n; j++ {
               if int64(2*a[i])+int64(a[j]) <= int64(m) && a[i]%a[j] == 0 {
                   adj[i] = append(adj[i], j)
               }
           }
       }
   }
   L := make([]int, n)
   R := make([]int, n)
   for i := 0; i < n; i++ {
       L[i] = -1
       R[i] = -1
   }
   vis := make([]bool, n)

   var dfs func(u int) bool
   dfs = func(u int) bool {
       if vis[u] {
           return false
       }
       vis[u] = true
       for _, v := range adj[u] {
           if L[v] == -1 || dfs(L[v]) {
               L[v] = u
               R[u] = v
               return true
           }
       }
       return false
   }
   // find matching
   for _, u := range big {
       for i := range vis {
           vis[i] = false
       }
       dfs(u)
   }
   // check unmatched big
   for _, u := range big {
       if R[u] == -1 {
           fmt.Fprintln(writer, -1)
           return
       }
   }
   // build answers
   type pair struct{ x, y int }
   ans := make([]pair, 0, n)
   for _, u := range big {
       v := R[u]
       ans = append(ans, pair{2*a[u] + a[v], a[u] + a[v]})
   }
   for i := 0; i < n; i++ {
       if a[i] <= m/3 && L[i] == -1 {
           ans = append(ans, pair{3 * a[i], 2 * a[i]})
       }
   }
   fmt.Fprintln(writer, len(ans))
   for _, p := range ans {
       fmt.Fprintf(writer, "%d %d\n", p.x, p.y)
   }
}
