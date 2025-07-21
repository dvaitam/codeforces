package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Edge represents a graph edge to node v with weight w
type Edge struct {
   to int
   w  int64
}

// Item is an entry in the priority queue
type Item struct {
   v    int
   dist int64
}

// PriorityQueue implements heap.Interface for Items
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

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   // build graph
   graph := make([][]Edge, n+1)
   roads := make([]struct{u, v int; w int64}, 0, m)
   for i := 0; i < m; i++ {
       var u, v int
       var w int64
       fmt.Fscan(reader, &u, &v, &w)
       graph[u] = append(graph[u], Edge{to: v, w: w})
       graph[v] = append(graph[v], Edge{to: u, w: w})
       roads = append(roads, struct{u, v int; w int64}{u, v, w})
   }
   trainsTo := make([][]int64, n+1)
   for i := 0; i < k; i++ {
       var s int
       var y int64
       fmt.Fscan(reader, &s, &y)
       trainsTo[s] = append(trainsTo[s], y)
       // add train edges
       graph[1] = append(graph[1], Edge{to: s, w: y})
       graph[s] = append(graph[s], Edge{to: 1, w: y})
   }
   // Dijkstra
   const INF = int64(4e18)
   dist := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = INF
   }
   dist[1] = 0
   pq := &PriorityQueue{{v: 1, dist: 0}}
   heap.Init(pq)
   for pq.Len() > 0 {
       it := heap.Pop(pq).(Item)
       u, d := it.v, it.dist
       if d != dist[u] {
           continue
       }
       for _, e := range graph[u] {
           v := e.to
           nd := d + e.w
           if nd < dist[v] {
               dist[v] = nd
               heap.Push(pq, Item{v: v, dist: nd})
           }
       }
   }
   // detect if each node has a road achieving shortest path
   hasRoad := make([]bool, n+1)
   for _, r := range roads {
       u, v, w := r.u, r.v, r.w
       if dist[u]+w == dist[v] {
           hasRoad[v] = true
       }
       if dist[v]+w == dist[u] {
           hasRoad[u] = true
       }
   }
   // count removable trains
   var ans int64
   for v := 2; v <= n; v++ {
       var match int64
       for _, y := range trainsTo[v] {
           if y > dist[v] {
               ans++
           } else if y == dist[v] {
               match++
           }
       }
       if match > 0 {
           if hasRoad[v] {
               ans += match
           } else {
               // keep one train, remove the rest
               ans += match - 1
           }
       }
   }
   fmt.Fprintln(writer, ans)
}
