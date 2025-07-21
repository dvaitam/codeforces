package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Edge represents an edge in the flow graph
type Edge struct {
   to, rev, cap int
   cost         int
}

// Graph is adjacency list of edges
type Graph [][]*Edge

// AddEdge adds edge from u to v with capacity cap and cost
func (g Graph) AddEdge(u, v, cap, cost int) {
   g[u] = append(g[u], &Edge{to: v, rev: len(g[v]), cap: cap, cost: cost})
   g[v] = append(g[v], &Edge{to: u, rev: len(g[u]) - 1, cap: 0, cost: -cost})
}

// Item for priority queue
type Item struct {
   v    int
   dist int
}

// Priority queue for Dijkstra
type PriorityQueue []Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
   old := *pq
   n := len(old)
   it := old[n-1]
   *pq = old[0 : n-1]
   return it
}

// minCostFlow returns (flow, cost). If can't send maxf flow, flow < maxf
func minCostFlow(g Graph, s, t, maxf int) (int, int) {
   n := len(g)
   const INF = 1<<60
   h := make([]int, n)      // potential
   prevv := make([]int, n)  // previous vertex
   preve := make([]int, n)  // previous edge index

   flow, cost := 0, 0
   for flow < maxf {
       dist := make([]int, n)
       for i := range dist {
           dist[i] = 1<<60
       }
       dist[s] = 0
       pq := &PriorityQueue{}
       heap.Init(pq)
       heap.Push(pq, Item{v: s, dist: 0})
       for pq.Len() > 0 {
           it := heap.Pop(pq).(Item)
           v := it.v
           if dist[v] < it.dist {
               continue
           }
           for i, e := range g[v] {
               if e.cap > 0 && dist[e.to] > dist[v]+e.cost+h[v]-h[e.to] {
                   dist[e.to] = dist[v] + e.cost + h[v] - h[e.to]
                   prevv[e.to] = v
                   preve[e.to] = i
                   heap.Push(pq, Item{v: e.to, dist: dist[e.to]})
               }
           }
       }
       if dist[t] == 1<<60 {
           break
       }
       for v := 0; v < n; v++ {
           if dist[v] < 1<<60 {
               h[v] += dist[v]
           }
       }
       d := maxf - flow
       // find minimum capacity along the path
       for v := t; v != s; v = prevv[v] {
           e := g[prevv[v]][preve[v]]
           if d > e.cap {
               d = e.cap
           }
       }
       flow += d
       cost += d * h[t]
       for v := t; v != s; v = prevv[v] {
           e := g[prevv[v]][preve[v]]
           e.cap -= d
           g[v][e.rev].cap += d
       }
   }
   return flow, cost
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t string
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   m := len(t)
   var n int
   fmt.Fscan(reader, &n)
   freqT := make([]int, 26)
   for i := 0; i < m; i++ {
       freqT[t[i]-'a']++
   }
   s := 0
   charBase := 1
   strBase := charBase + 26
   sink := strBase + n
   V := sink + 1
   g := make(Graph, V)
   // source to char nodes
   for c := 0; c < 26; c++ {
       if freqT[c] > 0 {
           g.AddEdge(s, charBase+c, freqT[c], 0)
       }
   }
   // string nodes
   for i := 1; i <= n; i++ {
       var si string
       var ai int
       fmt.Fscan(reader, &si, &ai)
       freqS := make([]int, 26)
       for j := 0; j < len(si); j++ {
           freqS[si[j]-'a']++
       }
       for c := 0; c < 26; c++ {
           if freqS[c] > 0 {
               // char node to string node with cost i
               g.AddEdge(charBase+c, strBase+(i-1), freqS[c], i)
           }
       }
       // string node to sink with cap ai
       if ai > 0 {
           g.AddEdge(strBase+(i-1), sink, ai, 0)
       }
   }
   if m == 0 {
       fmt.Fprintln(writer, 0)
       return
   }
   flow, cost := minCostFlow(g, s, sink, m)
   if flow < m {
       fmt.Fprintln(writer, -1)
   } else {
       fmt.Fprintln(writer, cost)
   }
}
