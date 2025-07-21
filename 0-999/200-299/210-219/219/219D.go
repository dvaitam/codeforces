package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge represents a connection to a neighboring city with a reversal cost
type Edge struct {
   to, cost int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   adj := make([][]Edge, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       // original directed edge u -> v has cost 0, reverse has cost 1
       adj[u] = append(adj[u], Edge{to: v, cost: 0})
       adj[v] = append(adj[v], Edge{to: u, cost: 1})
   }

   dp := make([]int, n+1)
   parent := make([]int, n+1)
   costToParent := make([]int, n+1)
   visited := make([]bool, n+1)
   queue := make([]int, 0, n)
   // BFS from 1 to compute initial reversals for capital = 1
   queue = append(queue, 1)
   visited[1] = true
   for i := 0; i < len(queue); i++ {
       u := queue[i]
       for _, e := range adj[u] {
           v, c := e.to, e.cost
           if !visited[v] {
               visited[v] = true
               dp[1] += c
               parent[v] = u
               costToParent[v] = c
               queue = append(queue, v)
           }
       }
   }

   // Reroot DP: compute dp[v] for all v
   for i := 1; i < len(queue); i++ {
       v := queue[i]
       u := parent[v]
       c := costToParent[v]
       dp[v] = dp[u] - c + (1 - c)
   }

   // Find minimum reversals and collect valid capitals
   minReversals := dp[1]
   for i := 2; i <= n; i++ {
       if dp[i] < minReversals {
           minReversals = dp[i]
       }
   }
   var result []int
   for i := 1; i <= n; i++ {
       if dp[i] == minReversals {
           result = append(result, i)
       }
   }

   fmt.Fprintln(out, minReversals)
   for idx, v := range result {
       if idx > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, v)
   }
   fmt.Fprintln(out)
}
