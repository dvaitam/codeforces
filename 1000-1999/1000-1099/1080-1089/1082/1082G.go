package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge represents an edge in the flow graph
type Edge struct {
   to   int
   rev  int
   cap  int64
}

// Dinic implements Dinic's max flow algorithm
type Dinic struct {
   n     int
   g     [][]Edge
   level []int
   iter  []int
}

// NewDinic creates a new Dinic for n nodes
func NewDinic(n int) *Dinic {
   g := make([][]Edge, n)
   level := make([]int, n)
   iter := make([]int, n)
   return &Dinic{n: n, g: g, level: level, iter: iter}
}

// AddEdge adds a directed edge from from->to with capacity cap
func (d *Dinic) AddEdge(from, to int, cap int64) {
   d.g[from] = append(d.g[from], Edge{to: to, rev: len(d.g[to]), cap: cap})
   d.g[to] = append(d.g[to], Edge{to: from, rev: len(d.g[from]) - 1, cap: 0})
}

// bfs builds level graph
func (d *Dinic) bfs(s int) {
   for i := range d.level {
       d.level[i] = -1
   }
   queue := make([]int, 0, d.n)
   d.level[s] = 0
   queue = append(queue, s)
   for qi := 0; qi < len(queue); qi++ {
       v := queue[qi]
       for _, e := range d.g[v] {
           if e.cap > 0 && d.level[e.to] < 0 {
               d.level[e.to] = d.level[v] + 1
               queue = append(queue, e.to)
           }
       }
   }
}

// dfs finds augmenting paths
func (d *Dinic) dfs(v, t int, f int64) int64 {
   if v == t {
       return f
   }
   for i := d.iter[v]; i < len(d.g[v]); i++ {
       e := &d.g[v][i]
       if e.cap > 0 && d.level[v] < d.level[e.to] {
           ret := d.dfs(e.to, t, min64(f, e.cap))
           if ret > 0 {
               e.cap -= ret
               d.g[e.to][e.rev].cap += ret
               return ret
           }
       }
       d.iter[v]++
   }
   return 0
}

// MaxFlow computes maximum flow from s to t
func (d *Dinic) MaxFlow(s, t int) int64 {
   var flow int64
   const INF = int64(1e18)
   for {
       d.bfs(s)
       if d.level[t] < 0 {
           break
       }
       for i := range d.iter {
           d.iter[i] = 0
       }
       for {
           f := d.dfs(s, t, INF)
           if f == 0 {
               break
           }
           flow += f
       }
   }
   return flow
}

func min64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   S := n + m
   T := S + 1
   d := NewDinic(T + 1)

   // capacities from nodes to sink
   for i := 0; i < n; i++ {
       var a int64
       fmt.Fscan(reader, &a)
       d.AddEdge(i, T, a)
   }

   var total int64
   // add item nodes and connect edges
   for i := 0; i < m; i++ {
       var u, v int
       var c int64
       fmt.Fscan(reader, &u, &v, &c)
       u--
       v--
       total += c
       item := n + i
       // S -> item
       d.AddEdge(S, item, c)
       // undirected connections item <-> u and item <-> v with cap c
       d.AddEdge(item, u, c)
       d.AddEdge(u, item, c)
       d.AddEdge(item, v, c)
       d.AddEdge(v, item, c)
   }

   maxf := d.MaxFlow(S, T)
   fmt.Fprintln(writer, total-maxf)
}
