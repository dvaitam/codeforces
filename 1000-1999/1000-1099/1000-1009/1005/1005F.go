package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge represents a graph edge with destination node and edge index
type Edge struct {
   to, idx int
}

var (
   n, m, k int
   adj      [][]Edge
   dis      []int
   s        []byte
   ans      int
   v        []string
)

// dfs recursively builds up to k shortest-path-parent selections for each node 2..n
func dfs(x int) {
   if ans >= k {
       return
   }
   if x == n+1 {
       ans++
       // record bits for edges 1..m
       v = append(v, string(s[1:m+1]))
       return
   }
   for _, e := range adj[x] {
       if dis[e.to]+1 == dis[x] {
           s[e.idx] = '1'
           dfs(x + 1)
           s[e.idx] = '0'
           if ans >= k {
               break
           }
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n, &m, &k)
   adj = make([][]Edge, n+1)
   for i := 1; i <= m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       adj[a] = append(adj[a], Edge{b, i})
       adj[b] = append(adj[b], Edge{a, i})
   }
   dis = make([]int, n+1)
   for i := 1; i <= n; i++ {
       dis[i] = -1
   }
   // BFS from node 1
   queue := make([]int, 0, n)
   dis[1] = 0
   queue = append(queue, 1)
   for i := 0; i < len(queue); i++ {
       x := queue[i]
       for _, e := range adj[x] {
           if dis[e.to] == -1 {
               dis[e.to] = dis[x] + 1
               queue = append(queue, e.to)
           }
       }
   }
   // initialize bit array for edges
   s = make([]byte, m+1)
   for i := 1; i <= m; i++ {
       s[i] = '0'
   }
   // generate up to k selections
   dfs(2)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
   for i := 0; i < ans; i++ {
       fmt.Fprintln(writer, v[i])
   }
}
