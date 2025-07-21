package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const INF = int64(4e18)

// Edge for min-cost max-flow
type Edge struct {
   to, rev, cap int
   cost, flow  int64
}

// Graph for MCMF
type MCMF struct {
   n       int
   graph   [][]Edge
   // potential []int64  // unused for SPFA
   dist    []int64
   prevV   []int
   prevE   []int
}

// NewMCMF creates a graph with n nodes
func NewMCMF(n int) *MCMF {
   return &MCMF{
       n: n,
       graph: make([][]Edge, n),
       // potential: make([]int64, n),
       dist: make([]int64, n),
       prevV: make([]int, n),
       prevE: make([]int, n),
   }
}

// AddEdge adds directed edge u->v
func (m *MCMF) AddEdge(u, v, cap int, cost int64) {
   m.graph[u] = append(m.graph[u], Edge{to: v, rev: len(m.graph[v]), cap: cap, cost: cost, flow: 0})
   m.graph[v] = append(m.graph[v], Edge{to: u, rev: len(m.graph[u]) - 1, cap: 0, cost: -cost, flow: 0})
}

// minCostFlow computes flow up to maxf, returns (flow, cost)
func (m *MCMF) minCostFlow(s, t, maxf int) (int, int64) {
   flow := 0
   cost := int64(0)
   for flow < maxf {
       // shortest path with SPFA
       inq := make([]bool, m.n)
       for i := 0; i < m.n; i++ {
           m.dist[i] = INF
       }
       m.dist[s] = 0
       queue := make([]int, 0, m.n)
       queue = append(queue, s)
       inq[s] = true
       head := 0
       for head < len(queue) {
           u := queue[head]
           head++
           inq[u] = false
           for ei, e := range m.graph[u] {
               if e.flow < int64(e.cap) && m.dist[u]+e.cost < m.dist[e.to] {
                   m.dist[e.to] = m.dist[u] + e.cost
                   m.prevV[e.to] = u
                   m.prevE[e.to] = ei
                   if !inq[e.to] {
                       inq[e.to] = true
                       queue = append(queue, e.to)
                   }
               }
           }
       }
       if m.dist[t] == INF {
           break
       }
       // add as much as possible (here 1 per path)
       df := maxf - flow
       // find bottleneck
       v := t
       for v != s {
           e := m.graph[m.prevV[v]][m.prevE[v]]
           if df > e.cap - int(e.flow) {
               df = e.cap - int(e.flow)
           }
           v = m.prevV[v]
       }
       if df <= 0 {
           break
       }
       flow += df
       // apply flow
       v = t
       for v != s {
           pe := &m.graph[m.prevV[v]][m.prevE[v]]
           pe.flow += int64(df)
           rev := &m.graph[v][pe.rev]
           rev.flow -= int64(df)
           cost += int64(df) * pe.cost
           v = m.prevV[v]
       }
   }
   return flow, cost
}


func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(in, &n, &k)
   s := make([]int64, n)
   t := make([]int64, n)
   c := make([]int64, n)
   times := make([]int64, 0, 2*n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &s[i], &t[i], &c[i])
       times = append(times, s[i])
       times = append(times, s[i]+t[i])
   }
   sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })
   uniq := times[:0]
   for i, v := range times {
       if i == 0 || v != times[i-1] {
           uniq = append(uniq, v)
       }
   }
   times = uniq
   m := len(times)
   // map time to index
   idx := func(x int64) int {
       i := sort.Search(len(times), func(i int) bool { return times[i] >= x })
       return i
   }
   // build graph
   V := m + 2
   S := m
   T := m + 1
   mf := NewMCMF(V)
   // source to first time
   mf.AddEdge(S, 0, k, 0)
   // last time to sink
   mf.AddEdge(m-1, T, k, 0)
   // timeline edges
   for i := 0; i < m-1; i++ {
       // capacity k, cost 0
       mf.AddEdge(i, i+1, k, 0)
   }
   // task edges
   taskEdge := make([]struct{u, ei int}, n)
   for i := 0; i < n; i++ {
       u := idx(s[i])
       v := idx(s[i] + t[i])
       uEdges := len(mf.graph[u])
       mf.AddEdge(u, v, 1, -c[i])
       taskEdge[i] = struct{u, ei int}{u, uEdges}
   }
   // flow
   mf.minCostFlow(S, T, k)
   // output
   res := make([]int, n)
   for i := 0; i < n; i++ {
       e := mf.graph[taskEdge[i].u][taskEdge[i].ei]
       if e.flow > 0 {
           res[i] = 1
       } else {
           res[i] = 0
       }
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i, v := range res {
       if i > 0 {
           w.WriteByte(' ')
       }
       if v == 1 {
           w.WriteByte('1')
       } else {
           w.WriteByte('0')
       }
   }
   w.WriteByte('\n')
}
