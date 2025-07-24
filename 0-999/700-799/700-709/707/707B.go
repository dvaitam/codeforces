package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

const INF = 1e18

type Edge struct {
   to   int
   w    int64
}

type Item struct {
   node int
   dist int64
}

// A PriorityQueue implements heap.Interface and holds Items.
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
   *pq = old[0 : n-1]
   return item
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   adj := make([][]Edge, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       var l int64
       fmt.Fscan(in, &u, &v, &l)
       adj[u] = append(adj[u], Edge{to: v, w: l})
       adj[v] = append(adj[v], Edge{to: u, w: l})
   }
   isStorage := make([]bool, n+1)
   storages := make([]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(in, &storages[i])
       isStorage[storages[i]] = true
   }
   // Multi-source Dijkstra
   dist := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = INF
   }
   pq := &PriorityQueue{}
   heap.Init(pq)
   for _, s := range storages {
       dist[s] = 0
       heap.Push(pq, Item{node: s, dist: 0})
   }
   for pq.Len() > 0 {
       cur := heap.Pop(pq).(Item)
       u, d := cur.node, cur.dist
       if d != dist[u] {
           continue
       }
       for _, e := range adj[u] {
           v := e.to
           nd := d + e.w
           if nd < dist[v] {
               dist[v] = nd
               heap.Push(pq, Item{node: v, dist: nd})
           }
       }
   }
   // Find minimal distance to non-storage
   ans := INF
   for i := 1; i <= n; i++ {
       if !isStorage[i] && dist[i] < ans {
           ans = dist[i]
       }
   }
   if ans == INF {
       fmt.Println(-1)
   } else {
       fmt.Println(ans)
   }
}
