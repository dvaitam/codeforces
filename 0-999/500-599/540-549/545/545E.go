package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Edge represents an adjacency list edge
type Edge struct {
   to int
   w  int64
   id int
}

// Item is an element in the priority queue
type Item struct {
   node int
   dist int64
}

// PriorityQueue implements heap.Interface for Items
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
   *pq = append(*pq, x.(*Item))
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
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   fmt.Fscan(in, &n, &m)
   graph := make([][]Edge, n+1)
   for i := 1; i <= m; i++ {
       var u, v int
       var w int64
       fmt.Fscan(in, &u, &v, &w)
       graph[u] = append(graph[u], Edge{to: v, w: w, id: i})
       graph[v] = append(graph[v], Edge{to: u, w: w, id: i})
   }
   var start int
   fmt.Fscan(in, &start)

   const INF = int64(1) << 62
   dis := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       dis[i] = INF
   }
   dis[start] = 0
   visited := make([]bool, n+1)

   pq := &PriorityQueue{}
   heap.Init(pq)
   heap.Push(pq, &Item{node: start, dist: 0})
   for pq.Len() > 0 {
       it := heap.Pop(pq).(*Item)
       u := it.node
       if visited[u] {
           continue
       }
       visited[u] = true
       for _, e := range graph[u] {
           if dis[e.to] > dis[u]+e.w {
               dis[e.to] = dis[u] + e.w
               heap.Push(pq, &Item{node: e.to, dist: dis[e.to]})
           }
       }
   }

   res := make([]int, 0, n-1)
   var total int64
   for u := 1; u <= n; u++ {
       if u == start {
           continue
       }
       minw := INF
       minid := 0
       for _, e := range graph[u] {
           if dis[u] == dis[e.to]+e.w && e.w < minw {
               minw = e.w
               minid = e.id
           }
       }
       total += minw
       res = append(res, minid)
   }

   fmt.Fprintln(out, total)
   for i, id := range res {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, id)
   }
   out.WriteByte('\n')
}
