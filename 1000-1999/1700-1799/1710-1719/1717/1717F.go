package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge represents a flow edge with capacity and reverse edge index.
type Edge struct {
   to, rev int
   cap     int64
}

// Graph is a flow graph implemented with adjacency lists.
type Graph struct {
   g [][]Edge
}

// NewGraph creates a new graph with n nodes.
func NewGraph(n int) *Graph {
   return &Graph{g: make([][]Edge, n)}
}

// AddEdge adds a directed edge from 'from' to 'to' with capacity cap.
// It returns the index of the edge in g[from] for later reference.
func (gr *Graph) AddEdge(from, to int, cap int64) int {
   idx := len(gr.g[from])
   gr.g[from] = append(gr.g[from], Edge{to: to, rev: len(gr.g[to]), cap: cap})
   gr.g[to] = append(gr.g[to], Edge{to: from, rev: idx, cap: 0})
   return idx
}

// min returns the smaller of a or b.
func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

// MaxFlow computes the maximum flow from s to t using Dinic's algorithm.
func (gr *Graph) MaxFlow(s, t int) int64 {
   n := len(gr.g)
   level := make([]int, n)
   iter := make([]int, n)
   const inf = int64(1e18)
   var bfs = func() {
       for i := range level {
           level[i] = -1
       }
       queue := make([]int, 0, n)
       level[s] = 0
       queue = append(queue, s)
       for qi := 0; qi < len(queue); qi++ {
           v := queue[qi]
           for _, e := range gr.g[v] {
               if e.cap > 0 && level[e.to] < 0 {
                   level[e.to] = level[v] + 1
                   queue = append(queue, e.to)
               }
           }
       }
   }
   var dfs func(v, t int, f int64) int64
   dfs = func(v, t int, f int64) int64 {
       if v == t {
           return f
       }
       for i := iter[v]; i < len(gr.g[v]); i++ {
           e := &gr.g[v][i]
           if e.cap > 0 && level[v] < level[e.to] {
               d := dfs(e.to, t, min(f, e.cap))
               if d > 0 {
                   e.cap -= d
                   gr.g[e.to][e.rev].cap += d
                   return d
               }
           }
           iter[v]++
       }
       return 0
   }
   flow := int64(0)
   for {
       bfs()
       if level[t] < 0 {
           break
       }
       for i := range iter {
           iter[i] = 0
       }
       for {
           f := dfs(s, t, inf)
           if f == 0 {
               break
           }
           flow += f
       }
   }
   return flow
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   s := make([]int64, n)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &s[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   totalNodes := n + m + 3
   S := n + m
   T := n + m + 1
   T2 := n + m + 2
   gr := NewGraph(totalNodes)
   cur := make([]int64, n)
   type pair struct{u,v,idxU,idxV int}
   edges := make([]pair, m)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       gr.AddEdge(S, i, 1)
       idxU := gr.AddEdge(i, m+u, 1)
       idxV := gr.AddEdge(i, m+v, 1)
       cur[u]--
       cur[v]--
       edges[i] = pair{u, v, idxU, idxV}
   }
   var ng bool
   tmp := int64(m)
   for i := 0; i < n; i++ {
       if s[i] == 0 {
           gr.AddEdge(m+i, T2, int64(1e18))
       } else if a[i] >= cur[i] {
           diff := a[i] - cur[i]
           if diff % 2 != 0 {
               ng = true
               break
           }
           cap := diff / 2
           gr.AddEdge(m+i, T, cap)
           tmp -= cap
       } else {
           ng = true
           break
       }
   }
   if tmp < 0 {
       ng = true
   }
   if ng {
       fmt.Fprintln(writer, "NO")
       return
   }
   gr.AddEdge(T2, T, tmp)
   res := gr.MaxFlow(S, T)
   if res != int64(m) {
       fmt.Fprintln(writer, "NO")
       return
   }
   fmt.Fprintln(writer, "YES")
   for i := 0; i < m; i++ {
       e := edges[i]
       // If edge to u was used (cap==0), orient v->u, else u->v
       if gr.g[e.u][e.idxU].cap == 0 {
           fmt.Fprintln(writer, e.v+1, e.u+1)
       } else {
           fmt.Fprintln(writer, e.u+1, e.v+1)
       }
   }
}
