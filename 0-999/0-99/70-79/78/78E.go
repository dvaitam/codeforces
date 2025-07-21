package main

import (
   "bufio"
   "fmt"
   "os"
)

// Dinic max flow implementation
type Edge struct {
   to, rev, cap int
}

type Dinic struct {
   n     int
   adj   [][]Edge
   level []int
   prog  []int
}

func NewDinic(n int) *Dinic {
   d := &Dinic{n: n, adj: make([][]Edge, n), level: make([]int, n), prog: make([]int, n)}
   return d
}

func (d *Dinic) AddEdge(u, v, c int) {
   d.adj[u] = append(d.adj[u], Edge{to: v, rev: len(d.adj[v]), cap: c})
   d.adj[v] = append(d.adj[v], Edge{to: u, rev: len(d.adj[u]) - 1, cap: 0})
}

func (d *Dinic) bfs(s, t int) bool {
   for i := range d.level {
       d.level[i] = -1
   }
   queue := make([]int, 0, d.n)
   d.level[s] = 0
   queue = append(queue, s)
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, e := range d.adj[u] {
           if e.cap > 0 && d.level[e.to] < 0 {
               d.level[e.to] = d.level[u] + 1
               queue = append(queue, e.to)
               if e.to == t {
                   return true
               }
           }
       }
   }
   return d.level[t] >= 0
}

func (d *Dinic) dfs(u, t, f int) int {
   if u == t {
       return f
   }
   for i := d.prog[u]; i < len(d.adj[u]); i++ {
       e := &d.adj[u][i]
       if e.cap > 0 && d.level[e.to] == d.level[u]+1 {
           minf := d.dfs(e.to, t, min(f, e.cap))
           if minf > 0 {
               e.cap -= minf
               d.adj[e.to][e.rev].cap += minf
               return minf
           }
       }
       d.prog[u]++
   }
   return 0
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func (d *Dinic) MaxFlow(s, t int) int {
   flow := 0
   for d.bfs(s, t) {
       for i := range d.prog {
           d.prog[i] = 0
       }
       for {
           f := d.dfs(s, t, 1<<30)
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
   var n, t int
   if _, err := fmt.Fscan(reader, &n, &t); err != nil {
       return
   }
   gridS := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &gridS[i])
   }
   // read blank line
   var blank string
   fmt.Fscan(reader, &blank)
   gridC := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &gridC[i])
   }
   // compute infection times
   const INF = 1e9
   dist := make([][]int, n)
   for i := range dist {
       dist[i] = make([]int, n)
       for j := range dist[i] {
           dist[i][j] = INF
       }
   }
   // find malfunctioning reactor Z position
   var zi, zj int
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if gridS[i][j] == 'Z' {
               zi, zj = i, j
           }
       }
   }
   // BFS
   type P struct{ i, j int }
   q := make([]P, 0, n*n)
   dist[zi][zj] = 0
   q = append(q, P{zi, zj})
   dirs := []P{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
   for qi := 0; qi < len(q); qi++ {
       u := q[qi]
       for _, ddir := range dirs {
           ni, nj := u.i+ddir.i, u.j+ddir.j
           if ni >= 0 && ni < n && nj >= 0 && nj < n && dist[ni][nj] == INF {
               // only labs (digits) are passable
               if gridS[ni][nj] >= '0' && gridS[ni][nj] <= '9' {
                   dist[ni][nj] = dist[u.i][u.j] + 1
                   q = append(q, P{ni, nj})
               }
           }
       }
   }
   // build flow network
   // node mapping
   timeLayers := t + 1
   totalLabs := n * n
   baseTime := 0
   baseCap := baseTime + totalLabs*timeLayers
   source := baseCap + totalLabs
   sink := source + 1
   V := sink + 1
   dinic := NewDinic(V)
   infCap := totalLabs * 10 + 5
   // helper functions
   timeNode := func(i, j, k int) int {
       return baseTime + (i*n+j)*timeLayers + k
   }
   capNode := func(i, j int) int {
       return baseCap + (i*n + j)
   }
   // add edges
   // source -> initial scientists at time 0
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if gridS[i][j] >= '0' && gridS[i][j] <= '9' {
               sCnt := int(gridS[i][j] - '0')
               if sCnt > 0 && dist[i][j] > 0 {
                   dinic.AddEdge(source, timeNode(i, j, 0), sCnt)
               }
           }
       }
   }
   // movement edges and capsule edges
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if gridS[i][j] >= '0' && gridS[i][j] <= '9' {
               // capsule count
               capCnt := 0
               if gridC[i][j] >= '0' && gridC[i][j] <= '9' {
                   capCnt = int(gridC[i][j] - '0')
               }
               // capNode->sink
               if capCnt > 0 {
                   dinic.AddEdge(capNode(i, j), sink, capCnt)
               }
               // for each time layer
               for k := 0; k <= t; k++ {
                   if dist[i][j] > k {
                       // capsule entry
                       if capCnt > 0 {
                           dinic.AddEdge(timeNode(i, j, k), capNode(i, j), infCap)
                       }
                       // movement to neighbors
                       if k < t {
                           for _, ddir := range dirs {
                               ni, nj := i+ddir.i, j+ddir.j
                               if ni >= 0 && ni < n && nj >= 0 && nj < n && dist[ni][nj] > k+1 && gridS[ni][nj] >= '0' && gridS[ni][nj] <= '9' {
                                   dinic.AddEdge(timeNode(i, j, k), timeNode(ni, nj, k+1), infCap)
                               }
                           }
                       }
                   }
               }
           }
       }
   }
   // compute max flow
   result := dinic.MaxFlow(source, sink)
   fmt.Println(result)
}
