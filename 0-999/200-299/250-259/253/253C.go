package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// Edge represents a weighted edge in graph
type Edge struct {
   to, w int
}

// Item for priority queue
type Item struct {
   v, dist int
}

// Priority queue of Items
type PriorityQueue []Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
   old := *pq
   n := len(old)
   it := old[n-1]
   *pq = old[:n-1]
   return it
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   var r1, c1, r2, c2 int
   fmt.Fscan(in, &r1, &c1, &r2, &c2)
   r1--
   r2--
   // compute row lengths +1
   b := make([]int, n)
   for i := 0; i < n; i++ {
       b[i] = a[i] + 1
   }
   // collect critical columns
   cols := make([]int, 0, n+2)
   cols = append(cols, c1, c2)
   for i := 0; i < n; i++ {
       cols = append(cols, b[i])
   }
   sort.Ints(cols)
   // unique
   uniq := cols[:0]
   for i, v := range cols {
       if i == 0 || v != cols[i-1] {
           uniq = append(uniq, v)
       }
   }
   cols = uniq
   m := len(cols)
   // map col value to index
   colIndex := make(map[int]int, m)
   for i, v := range cols {
       colIndex[v] = i
   }
   // valid positions per row
   valid := make([][]int, n)
   for i := 0; i < n; i++ {
       // cols <= b[i]
       for j, v := range cols {
           if v <= b[i] {
               valid[i] = append(valid[i], j)
           } else {
               break
           }
       }
   }
   // graph adjacency
   N := n * m
   adj := make([][]Edge, N)
   // helper to node id
   id := func(r, ci int) int { return r*m + ci }
   // horizontal edges
   for i := 0; i < n; i++ {
       vs := valid[i]
       for k := 1; k < len(vs); k++ {
           u := vs[k-1]
           v := vs[k]
           w := cols[v] - cols[u]
           ui := id(i, u)
           vi := id(i, v)
           adj[ui] = append(adj[ui], Edge{vi, w})
           adj[vi] = append(adj[vi], Edge{ui, w})
       }
   }
   // vertical edges
   for i := 0; i < n; i++ {
       for _, u := range valid[i] {
           cu := cols[u]
           ui := id(i, u)
           // up
           if i > 0 {
               // target col
               tc := cu
               if tc > b[i-1] {
                   tc = b[i-1]
               }
               vi, ok := colIndex[tc]
               if ok {
                   adj[ui] = append(adj[ui], Edge{id(i-1, vi), 1})
               }
           }
           // down
           if i+1 < n {
               tc := cu
               if tc > b[i+1] {
                   tc = b[i+1]
               }
               vi, ok := colIndex[tc]
               if ok {
                   adj[ui] = append(adj[ui], Edge{id(i+1, vi), 1})
               }
           }
       }
   }
   // Dijkstra
   const INF = 1e18
   dist := make([]int, N)
   for i := range dist {
       dist[i] = 1e18
   }
   // start and target ids
   si := id(r1, colIndex[c1])
   ti := id(r2, colIndex[c2])
   dist[si] = 0
   pq := &PriorityQueue{{si, 0}}
   heap.Init(pq)
   for pq.Len() > 0 {
       it := heap.Pop(pq).(Item)
       v, d := it.v, it.dist
       if d != dist[v] {
           continue
       }
       if v == ti {
           break
       }
       for _, e := range adj[v] {
           nd := d + e.w
           if nd < dist[e.to] {
               dist[e.to] = nd
               heap.Push(pq, Item{e.to, nd})
           }
       }
   }
   fmt.Println(dist[ti])
}
