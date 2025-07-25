package main

import (
   "container/heap"
   "bufio"
   "fmt"
   "os"
)

// Edge represents an edge in the flow graph
type Edge struct {
   to, rev, cap, cost int
}

// Graph is adjacency list
var graph [][]Edge

func addEdge(u, v, cap, cost int) {
   graph[u] = append(graph[u], Edge{v, len(graph[v]), cap, cost})
   graph[v] = append(graph[v], Edge{u, len(graph[u]) - 1, 0, -cost})
}

// Item for priority queue
type Item struct {
   v, dist int
}
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
   in := bufio.NewReader(os.Stdin)
   var n, m, k, c, d int
   fmt.Fscan(in, &n, &m, &k, &c, &d)
   a := make([]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(in, &a[i])
       a[i]--
   }
   adj := make([][]int, n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--, v--
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // BFS to compute dist to node 0
   const INF = 1e9
   dist0 := make([]int, n)
   for i := range dist0 { dist0[i] = INF }
   queue := make([]int, 0, n)
   dist0[0] = 0
   queue = append(queue, 0)
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, v := range adj[u] {
           if dist0[v] > dist0[u]+1 {
               dist0[v] = dist0[u] + 1
               queue = append(queue, v)
           }
       }
   }
   maxd := 0
   for _, x := range a {
       if dist0[x] > maxd {
           maxd = dist0[x]
       }
   }
   T := maxd + k
   // build time-expanded graph
   numNodes := 2 + (T+1)*n
   S := 0
   sink := numNodes - 1
   graph = make([][]Edge, numNodes)
   // source to each start
   for _, x := range a {
       u := 1 + x
       addEdge(S, u, 1, 0)
   }
   // time layers
   for t := 0; t < T; t++ {
       for i := 0; i < n; i++ {
           u := 1 + t*n + i
           v := 1 + (t+1)*n + i
           // wait
           addEdge(u, v, k, 0)
       }
       // move along edges
       for u0 := 0; u0 < n; u0++ {
           for _, v0 := range adj[u0] {
               u := 1 + t*n + u0
               v := 1 + (t+1)*n + v0
               // add k unit arcs with marginal costs
               for x := 1; x <= k; x++ {
                   cost := d * (2*x - 1)
                   addEdge(u, v, 1, cost)
               }
           }
       }
   }
   // sink edges from node 0 at each time
   for t := 0; t <= T; t++ {
       u := 1 + t*n + 0
       addEdge(u, sink, k, c*t)
   }
   // min cost max flow
   N := numNodes
   flow := 0
   costRes := 0
   potential := make([]int, N)
   dist := make([]int, N)
   prevv := make([]int, N)
   preve := make([]int, N)
   // successive shortest path
   for flow < k {
       // dijkstra
       for i := 0; i < N; i++ { dist[i] = INF }
       dist[S] = 0
       hq := &PriorityQueue{{S, 0}}
       heap.Init(hq)
       for hq.Len() > 0 {
           it := heap.Pop(hq).(Item)
           v := it.v
           if dist[v] < it.dist {
               continue
           }
           for ei, e := range graph[v] {
               if e.cap > 0 && dist[e.to] > dist[v]+e.cost+potential[v]-potential[e.to] {
                   dist[e.to] = dist[v] + e.cost + potential[v] - potential[e.to]
                   prevv[e.to] = v
                   preve[e.to] = ei
                   heap.Push(hq, Item{e.to, dist[e.to]})
               }
           }
       }
       if dist[sink] == INF {
           break
       }
       for v := 0; v < N; v++ {
           if dist[v] < INF {
               potential[v] += dist[v]
           }
       }
       // add one flow
       addf := 1
       flow += addf
       v := sink
       for v != S {
           e := &graph[prevv[v]][preve[v]]
           e.cap -= addf
           graph[v][e.rev].cap += addf
           v = prevv[v]
       }
       costRes += potential[sink] * addf
   }
   fmt.Println(costRes)
}
