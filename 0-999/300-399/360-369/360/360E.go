package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "io"
   "os"
)

type Edge struct {
   to, w, id int
}

// Item for priority queue
type Item struct {
   v   int
   dist int64
}

type PQ []Item
func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PQ) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
   old := *pq
   n := len(old)
   it := old[n-1]
   *pq = old[0 : n-1]
   return it
}


func dijkstraDist(n int, adj [][]Edge, s int) []int64 {
   const INF = 1<<62
   dist := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = INF
   }
   dist[s] = 0
   pq := &PQ{}
   heap.Init(pq)
   heap.Push(pq, Item{s, 0})
   for pq.Len() > 0 {
       it := heap.Pop(pq).(Item)
       v, d := it.v, it.dist
       if d != dist[v] {
           continue
       }
       for _, e := range adj[v] {
           nd := d + int64(e.w)
           if nd < dist[e.to] {
               dist[e.to] = nd
               heap.Push(pq, Item{e.to, nd})
           }
       }
   }
   return dist
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   var s1, s2, f int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       if err == io.EOF { return }
       panic(err)
   }
   fmt.Fscan(in, &s1, &s2, &f)
   fixed := make([][3]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &fixed[i][0], &fixed[i][1], &fixed[i][2])
   }
   au := make([]int, k)
   av := make([]int, k)
   al := make([]int, k)
   ar := make([]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(in, &au[i], &av[i], &al[i], &ar[i])
   }
   // Simple strategy: set all adjustable edges to minimum
   x := make([]int, k)
   for i := 0; i < k; i++ {
       x[i] = al[i]
   }
   // Build graph with weights x
   adj := make([][]Edge, n+1)
   for _, ed := range fixed {
       adj[ed[0]] = append(adj[ed[0]], Edge{ed[1], ed[2], -1})
   }
   for i := 0; i < k; i++ {
       adj[au[i]] = append(adj[au[i]], Edge{av[i], x[i], i})
   }
   d1 := dijkstraDist(n, adj, s1)[f]
   d2 := dijkstraDist(n, adj, s2)[f]
   if d1 < d2 {
       fmt.Println("WIN")
       for i := 0; i < k; i++ {
           fmt.Printf("%d ", x[i])
       }
       fmt.Println()
   } else if d1 == d2 {
       fmt.Println("DRAW")
       for i := 0; i < k; i++ {
           fmt.Printf("%d ", x[i])
       }
       fmt.Println()
   } else {
       fmt.Println("LOSE")
   }
}
