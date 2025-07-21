package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Edge represents an edge in the graph
type Edge struct {
   to     int
   weight int64
}

// Item is a vertex in the priority queue
type Item struct {
   vertex int
   dist   int64
   index  int // index in the heap
}

// PriorityQueue implements heap.Interface for Items
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
   return pq[i].dist < pq[j].dist
}
func (pq PriorityQueue) Swap(i, j int) {
   pq[i], pq[j] = pq[j], pq[i]
   pq[i].index = i
   pq[j].index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
   n := len(*pq)
   item := x.(*Item)
   item.index = n
   *pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
   old := *pq
   n := len(old)
   item := old[n-1]
   item.index = -1 // for safety
   *pq = old[0 : n-1]
   return item
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   adj := make([][]Edge, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       var w int64
       fmt.Fscan(reader, &u, &v, &w)
       adj[u] = append(adj[u], Edge{to: v, weight: w})
       adj[v] = append(adj[v], Edge{to: u, weight: w})
   }
   const inf = int64(1e18)
   dist := make([]int64, n+1)
   prev := make([]int, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = inf
       prev[i] = 0
   }
   dist[1] = 0

   pq := &PriorityQueue{}
   heap.Init(pq)
   heap.Push(pq, &Item{vertex: 1, dist: 0})

   for pq.Len() > 0 {
       item := heap.Pop(pq).(*Item)
       u := item.vertex
       d := item.dist
       if d != dist[u] {
           continue
       }
       if u == n {
           break
       }
       for _, e := range adj[u] {
           v := e.to
           nd := d + e.weight
           if nd < dist[v] {
               dist[v] = nd
               prev[v] = u
               heap.Push(pq, &Item{vertex: v, dist: nd})
           }
       }
   }
   if dist[n] == inf {
       fmt.Fprintln(writer, -1)
       return
   }
   // reconstruct path
   path := make([]int, 0)
   for v := n; v != 0; v = prev[v] {
       path = append(path, v)
   }
   // reverse
   for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
       path[i], path[j] = path[j], path[i]
   }
   // output path
   for i, v := range path {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprintf(writer, "%d", v)
   }
   writer.WriteByte('\n')
}
