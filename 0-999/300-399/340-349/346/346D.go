package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

const INF = 1000000007

type Item struct {
   v, f int
}
type PQ []Item
func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool { return pq[i].f < pq[j].f }
func (pq PQ) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
   old := *pq
   n := len(old)
   x := old[n-1]
   *pq = old[0 : n-1]
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   fwd := make([][]int, n+1)
   rev := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       fwd[u] = append(fwd[u], v)
       rev[v] = append(rev[v], u)
   }
   var s, t int
   fmt.Fscan(in, &s, &t)
   // BFS from s
   reachS := make([]bool, n+1)
   queue := make([]int, 0, n)
   reachS[s] = true
   queue = append(queue, s)
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, v := range fwd[u] {
           if !reachS[v] {
               reachS[v] = true
               queue = append(queue, v)
           }
       }
   }
   // BFS from t on reverse
   reachT := make([]bool, n+1)
   queue = queue[:0]
   reachT[t] = true
   queue = append(queue, t)
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, v := range rev[u] {
           if !reachT[v] {
               reachT[v] = true
               queue = append(queue, v)
           }
       }
   }
   if !reachS[t] {
       fmt.Fprintln(out, -1)
       return
   }
   good := make([]bool, n+1)
   for i := 1; i <= n; i++ {
       good[i] = reachS[i] && reachT[i]
   }
   deg := make([]int, n+1)
   revGood := make([][]int, n+1)
   for u := 1; u <= n; u++ {
       if !good[u] {
           continue
       }
       for _, v := range fwd[u] {
           if good[v] {
               deg[u]++
               revGood[v] = append(revGood[v], u)
           }
       }
   }
   // DP via Dijkstra-like
   f := make([]int, n+1)
   minF := make([]int, n+1)
   maxF := make([]int, n+1)
   cnt := make([]int, n+1)
   for i := 1; i <= n; i++ {
       f[i] = INF
       minF[i] = INF
   }
   f[t] = 0
   // priority queue
   pq := &PQ{}
   heap.Init(pq)
   heap.Push(pq, Item{v: t, f: 0})
   for pq.Len() > 0 {
       it := heap.Pop(pq).(Item)
       v := it.v; fv := it.f
       if fv != f[v] {
           continue
       }
       for _, u := range revGood[v] {
           cnt[u]++
           if fv < minF[u] {
               minF[u] = fv
           }
           if fv > maxF[u] {
               maxF[u] = fv
           }
           if cnt[u] == deg[u] {
               // compute f[u]
               noOrder := maxF[u]
               order := minF[u] + 1
               cand := noOrder
               if order < cand {
                   cand = order
               }
               if cand < f[u] {
                   f[u] = cand
                   heap.Push(pq, Item{v: u, f: cand})
               }
           }
       }
   }
   if f[s] >= INF {
       fmt.Fprintln(out, -1)
   } else {
       fmt.Fprintln(out, f[s])
   }
}
