package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "math"
   "os"
   "sort"
   "math/bits"
)

const INF = math.MaxInt64 / 4

// Edge represents a graph edge with weight
type Edge struct {
   to int
   w  int64
}

// Item is a node in the priority queue
type Item struct {
   node int
   dist int64
}

// Priority queue implementation
type PriorityQueue []Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
   old := *pq
   n := len(old)
   item := old[n-1]
   *pq = old[0 : n-1]
   return item
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   // read graph
   adj := make([][]Edge, n)
   type EdgeRec struct{ u, v int; w int64 }
   edges := make([]EdgeRec, m)
   for i := 0; i < m; i++ {
       var x, y int
       var l int64
       fmt.Fscan(reader, &x, &y, &l)
       x--; y--
       adj[x] = append(adj[x], Edge{y, l})
       adj[y] = append(adj[y], Edge{x, l})
       edges[i] = EdgeRec{x, y, l}
   }
   // prepare result buffer
   totalPairs := n * (n - 1) / 2
   result := make([]int, 0, totalPairs)
   // parameters for bitsets
   blocks := (n + 63) / 64

   // helper slices
   dist := make([]int64, n)
   inDeg := make([]int, n)
   // DAG adjacency
   var adjDAG [][]int
   // bitsets
   Reach := make([][]uint64, n)
   for i := 0; i < n; i++ {
       Reach[i] = make([]uint64, blocks)
   }

   // process each source
   for s := 0; s < n; s++ {
       // Dijkstra
       for i := 0; i < n; i++ {
           dist[i] = INF
       }
       dist[s] = 0
       pq := &PriorityQueue{}
       heap.Init(pq)
       heap.Push(pq, Item{s, 0})
       for pq.Len() > 0 {
           it := heap.Pop(pq).(Item)
           u, d := it.node, it.dist
           if d != dist[u] {
               continue
           }
           for _, e := range adj[u] {
               nd := d + e.w
               if nd < dist[e.to] {
                   dist[e.to] = nd
                   heap.Push(pq, Item{e.to, nd})
               }
           }
       }
       // build shortest-path DAG
       adjDAG = make([][]int, n)
       for i := range inDeg {
           inDeg[i] = 0
       }
       for _, e := range edges {
           u, v, w := e.u, e.v, e.w
           if dist[u] + w == dist[v] {
               adjDAG[u] = append(adjDAG[u], v)
               inDeg[v]++
           }
           if dist[v] + w == dist[u] {
               adjDAG[v] = append(adjDAG[v], u)
               inDeg[u]++
           }
       }
       // initialize reachability bitsets
       for i := 0; i < n; i++ {
           for b := 0; b < blocks; b++ {
               Reach[i][b] = 0
           }
           Reach[i][i/64] |= 1 << (uint(i) % 64)
       }
       // nodes sorted by descending dist
       nodes := make([]int, n)
       for i := 0; i < n; i++ {
           nodes[i] = i
       }
       sort.Slice(nodes, func(i, j int) bool { return dist[nodes[i]] > dist[nodes[j]] })
       // propagate reach sets
       for _, v := range nodes {
           for _, w := range adjDAG[v] {
               // OR Reach[w] into Reach[v]
               for b := 0; b < blocks; b++ {
                   Reach[v][b] |= Reach[w][b]
               }
           }
       }
       // compute answers for this source
       ans := make([]int, n)
       for v := 0; v < n; v++ {
           deg := inDeg[v]
           if deg == 0 {
               continue
           }
           for t := s + 1; t < n; t++ {
               if (Reach[v][t/64]>>(uint(t)%64))&1 == 1 {
                   ans[t] += deg
               }
           }
       }
       // collect results for pairs s<t
       for t := s + 1; t < n; t++ {
           result = append(result, ans[t])
       }
   }
   // output
   for i, v := range result {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}
