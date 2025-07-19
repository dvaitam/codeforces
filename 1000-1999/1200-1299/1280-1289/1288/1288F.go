package main

import (
   "bufio"
   "fmt"
   "os"
)

type Edge struct {
   to, rev int
   cap, cost int
}

func min(a, b int) int {
   if a < b { return a }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n1, n2, m, R, B int
   fmt.Fscan(reader, &n1, &n2, &m, &R, &B)
   var s1, s2 string
   fmt.Fscan(reader, &s1)
   fmt.Fscan(reader, &s2)
   N := n1 + n2 + 4
   s := n1 + n2
   t := s + 1
   S := t + 1
   T := S + 1
   graph := make([][]Edge, N)
   deg := make([]int, N)
   // function to add edge
   addEdge := func(u, v, cap, cost int) {
       graph[u] = append(graph[u], Edge{to: v, rev: len(graph[v]), cap: cap, cost: cost})
       graph[v] = append(graph[v], Edge{to: u, rev: len(graph[u]) - 1, cap: 0, cost: -cost})
   }
   // add edge with lower bound l and upper bound r
   addLower := func(u, v, l, r, cost int) {
       deg[u] -= l
       deg[v] += l
       addEdge(u, v, r - l, cost)
   }
   INF_CAP := m + n1 + n2 + 5
   // source-sink coloring edges
   for i := 0; i < n1; i++ {
       c := s1[i]
       if c == 'R' {
           addLower(s, i, 1, INF_CAP, 0)
       } else if c == 'B' {
           addLower(i, t, 1, INF_CAP, 0)
       } else {
           addLower(s, i, 0, INF_CAP, 0)
           addLower(i, t, 0, INF_CAP, 0)
       }
   }
   for i := 0; i < n2; i++ {
       c := s2[i]
       node := n1 + i
       if c == 'R' {
           addLower(node, t, 1, INF_CAP, 0)
       } else if c == 'B' {
           addLower(s, node, 1, INF_CAP, 0)
       } else {
           addLower(s, node, 0, INF_CAP, 0)
           addLower(node, t, 0, INF_CAP, 0)
       }
   }
   // record m edges indices
   us := make([]int, m)
   vs := make([]int, m)
   reIdx := make([]int, m)
   beIdx := make([]int, m)
   // add edges between groups
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       v = v + n1
       us[i] = u
       vs[i] = v
       // R color
       deg[u] -= 0; deg[v] += 0
       addEdge(u, v, 1, R)
       reIdx[i] = len(graph[u]) - 1
       // B color
       addEdge(v, u, 1, B)
       beIdx[i] = len(graph[v]) - 1
   }
   // add circulation edge t->s
   addEdge(t, s, INF_CAP, 0)
   // super source/sink for lower bounds
   sum := 0
   for i := 0; i < N; i++ {
       if deg[i] > 0 {
           addEdge(S, i, deg[i], 0)
           sum += deg[i]
       } else if deg[i] < 0 {
           addEdge(i, T, -deg[i], 0)
       }
   }
   // min cost flow from S to T
   const INF_COST = 1<<60
   flow, cost := 0, 0
   dist := make([]int, N)
   inQ := make([]bool, N)

   for {
       // SPFA
       for i := 0; i < N; i++ {
           dist[i] = INF_COST
           inQ[i] = false
       }
       queue := make([]int, 0, N)
       dist[S] = 0
       queue = append(queue, S)
       inQ[S] = true
       for qi := 0; qi < len(queue); qi++ {
           u := queue[qi]
           inQ[u] = false
           for _, e := range graph[u] {
               if e.cap > 0 && dist[e.to] > dist[u] + e.cost {
                   dist[e.to] = dist[u] + e.cost
                   if !inQ[e.to] {
                       inQ[e.to] = true
                       queue = append(queue, e.to)
                   }
               }
           }
       }
       if dist[T] == INF_COST {
           break
       }
       // find blocking flow
       cur := make([]int, N)
       var dfs func(int, int) int
       visited := make([]bool, N)
       dfs = func(u, f int) int {
           if u == T {
               flow += f
               cost += f * dist[T]
               return f
           }
           visited[u] = true
           for i := cur[u]; i < len(graph[u]); i++ {
               e := &graph[u][i]
               if e.cap > 0 && !visited[e.to] && dist[e.to] == dist[u] + e.cost {
                   pushed := dfs(e.to, min(f, e.cap))
                   if pushed > 0 {
                       e.cap -= pushed
                       graph[e.to][e.rev].cap += pushed
                       return pushed
                   }
               }
               cur[u]++
           }
           return 0
       }
       for {
           for i := range visited { visited[i] = false }
           pushed := dfs(S, INF_CAP)
           if pushed == 0 {
               break
           }
       }
   }
   if flow < sum {
       fmt.Fprintln(writer, -1)
       return
   }
   // output cost and assignment
   fmt.Fprintln(writer, cost)
   res := make([]byte, m)
   for i := 0; i < m; i++ {
       if graph[us[i]][reIdx[i]].cap == 0 {
           res[i] = 'R'
       } else if graph[vs[i]][beIdx[i]].cap == 0 {
           res[i] = 'B'
       } else {
           res[i] = 'U'
       }
   }
   writer.Write(res)
   writer.WriteByte('\n')
}
