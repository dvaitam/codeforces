package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Edge represents an edge to a node with a certain weight
type Edge struct {
   to   int
   cost int64
}

// Item is a node in the priority queue
type Item struct {
   node int
   dist int64
}

// PriorityQueue implements heap.Interface for Items
type PriorityQueue []Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) {
   *pq = append(*pq, x.(Item))
}
func (pq *PriorityQueue) Pop() interface{} {
   old := *pq
   n := len(old)
   item := old[n-1]
   *pq = old[:n-1]
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
       adj[u] = append(adj[u], Edge{to: v, cost: w})
       adj[v] = append(adj[v], Edge{to: u, cost: w})
   }

   const INF = int64(1e18)
   dist := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = INF
   }
   dist[1] = 0

   pq := &PriorityQueue{}
   heap.Init(pq)
   heap.Push(pq, Item{node: 1, dist: 0})

   for pq.Len() > 0 {
       cur := heap.Pop(pq).(Item)
       if cur.dist != dist[cur.node] {
           continue
       }
       if cur.node == n {
           break
       }
       for _, e := range adj[cur.node] {
           nd := cur.dist + e.cost
           if nd < dist[e.to] {
               dist[e.to] = nd
               heap.Push(pq, Item{node: e.to, dist: nd})
           }
       }
   }

   if dist[n] == INF {
       fmt.Fprintln(writer, -1)
   } else {
       fmt.Fprintln(writer, dist[n])
   }
}
