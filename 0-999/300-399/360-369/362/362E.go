package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

const INF = int64(1e18)

// Edge represents an edge in residual graph
type Edge struct {
   to, rev int
   cap, cost int64
}

// Graph is adjacency list of edges
var graph [][]Edge
var h, dist []int64
var prevv, preve []int

func addEdge(from, to int, cap, cost int64) {
   graph[from] = append(graph[from], Edge{to, len(graph[to]), cap, cost})
   graph[to] = append(graph[to], Edge{from, len(graph[from]) - 1, 0, -cost})
}

// Item for priority queue
type Item struct {
   v    int
   dist int64
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
   rdr := bufio.NewReader(os.Stdin)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   var n int
   var k int64
   fmt.Fscan(rdr, &n, &k)
   c := make([][]int64, n)
   for i := 0; i < n; i++ {
       c[i] = make([]int64, n)
       for j := 0; j < n; j++ {
           fmt.Fscan(rdr, &c[i][j])
       }
   }
   // build graph
   graph = make([][]Edge, n)
   s, t := 0, n-1
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if c[i][j] > 0 {
               addEdge(i, j, c[i][j], 0)
               addEdge(i, j, k, 1)
           }
       }
   }
   // prepare mcmf
   h = make([]int64, n)
   dist = make([]int64, n)
   prevv = make([]int, n)
   preve = make([]int, n)
   flow := int64(0)
   costSpent := int64(0)
   // successive shortest augmenting path
   for {
       // dijkstra
       for i := 0; i < n; i++ {
           dist[i] = INF
       }
       dist[s] = 0
       pq := &PriorityQueue{{s, 0}}
       heap.Init(pq)
       for pq.Len() > 0 {
           it := heap.Pop(pq).(Item)
           v := it.v
           if dist[v] < it.dist {
               continue
           }
           for ei, e := range graph[v] {
               if e.cap > 0 && dist[e.to] > dist[v] + e.cost + h[v] - h[e.to] {
                   dist[e.to] = dist[v] + e.cost + h[v] - h[e.to]
                   prevv[e.to] = v
                   preve[e.to] = ei
                   heap.Push(pq, Item{e.to, dist[e.to]})
               }
           }
       }
       if dist[t] == INF {
           break
       }
       for v := 0; v < n; v++ {
           if dist[v] < INF {
               h[v] += dist[v]
           }
       }
       // determine amount to send
       d := INF
       v := t
       for v != s {
           e := graph[prevv[v]][preve[v]]
           if d > e.cap {
               d = e.cap
           }
           v = prevv[v]
       }
       costPerUnit := h[t]
       var amt int64
       if costPerUnit > 0 {
           maxAdd := (k - costSpent) / costPerUnit
           if maxAdd <= 0 {
               break
           }
           if d > maxAdd {
               d = maxAdd
           }
           amt = d
       } else {
           amt = d
       }
       // apply flow
       v = t
       for v != s {
           e := &graph[prevv[v]][preve[v]]
           e.cap -= amt
           // reverse edge
           graph[v][e.rev].cap += amt
           v = prevv[v]
       }
       flow += amt
       costSpent += amt * costPerUnit
   }
   fmt.Fprintln(w, flow)
}
