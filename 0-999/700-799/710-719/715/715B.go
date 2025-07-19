package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

const INF int64 = 1 << 60

// edge represents a directed edge in graph
type edge struct {
   to, w, nxt int
}

// Item for priority queue
type Item struct {
   dist int64
   node int
}

// PriorityQueue implements heap.Interface
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

func dijkstra(start int, d []int64, first []int, e []edge) {
   for i := range d {
      d[i] = INF
   }
   d[start] = 0
   pq := &PriorityQueue{}
   heap.Init(pq)
   heap.Push(pq, Item{0, start})
   for pq.Len() > 0 {
      u := heap.Pop(pq).(Item)
      if d[u.node] != u.dist {
         continue
      }
      for i := first[u.node]; i != -1; i = e[i].nxt {
         v := e[i].to
         nd := u.dist + int64(e[i].w)
         if d[v] > nd {
            d[v] = nd
            heap.Push(pq, Item{nd, v})
         }
      }
   }
}

func findx(s, t int, first []int, e []edge, d []int64, pre []int, vis []bool) {
   for i := range vis {
      vis[i] = false
   }
   queue := make([]int, 0, len(first))
   queue = append(queue, s)
   vis[s] = true
   for qi := 0; qi < len(queue); qi++ {
      u := queue[qi]
      if u == t {
         return
      }
      for i := first[u]; i != -1; i = e[i].nxt {
         v := e[i].to
         if vis[v] {
            continue
         }
         if d[v] >= d[u]+int64(e[i].w) {
            pre[v] = i ^ 1
            if v == t {
               return
            }
            vis[v] = true
            queue = append(queue, v)
         }
      }
   }
}

func max64(a, b int64) int64 {
   if a > b {
      return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m int
   var L int64
   var s, t int
   fmt.Fscan(in, &n, &m, &L, &s, &t)
   first := make([]int, n)
   for i := range first {
      first[i] = -1
   }
   e := make([]edge, 0, 2*m)
   cho := make([]bool, 2*m)
   pre := make([]int, n)
   var ec int

   // add edge u->v with weight w
   addEdge := func(u, v, w int) {
      e = append(e, edge{to: v, w: w, nxt: first[u]})
      first[u] = ec
      ec++
   }

   for i := 0; i < m; i++ {
      var u, v, w int
      fmt.Fscan(in, &u, &v, &w)
      addEdge(u, v, w)
      addEdge(v, u, w)
   }
   // initialize zero edges to INF, others mark cho
   for i := 0; i < ec; i++ {
      if e[i].w == 0 {
         e[i].w = int(INF)
      } else {
         cho[i] = true
      }
   }
   d2 := make([]int64, n)
   dijkstra(t, d2, first, e)
   if d2[s] < L {
      fmt.Fprint(out, "NO")
      return
   }
   // restore zero-weight edges to weight 1
   for i := 0; i < ec; i++ {
      if e[i].w > 1000000000 {
         e[i].w = 1
      }
   }
   d1 := make([]int64, n)
   dijkstra(s, d1, first, e)
   dt := make([]int64, n)
   for i := 0; i < n; i++ {
      dt[i] = max64(d1[i], L-d2[i])
   }
   d := make([]int64, n)
   for i := range d {
      d[i] = INF
   }
   d[s] = 0
   pq := &PriorityQueue{}
   heap.Init(pq)
   heap.Push(pq, Item{0, s})
   for pq.Len() > 0 {
      u := heap.Pop(pq).(Item)
      if d[u.node] != u.dist {
         continue
      }
      for i := first[u.node]; i != -1; i = e[i].nxt {
         v := e[i].to
         nd := u.dist + int64(e[i].w)
         tval := max64(dt[v], nd)
         if d[v] > tval {
            d[v] = tval
            heap.Push(pq, Item{tval, v})
         }
      }
   }
   vis := make([]bool, n)
   flag := false
   for i := 0; i < n; i++ {
      if d[i]+d2[i] <= L {
         // adjust path
         findx(s, i, first, e, d, pre, vis)
         u := i
         now := L - d2[i]
         for u != s {
            x := pre[u]
            v := e[x].to
            if !cho[x] {
               cho[x] = true
               e[x].w = int(now - d[v])
               now = d[v]
            } else {
               now -= int64(e[x].w)
            }
            u = v
         }
         flag = true
         break
      }
   }
   if !flag {
      fmt.Fprint(out, "NO")
      return
   }
   fmt.Fprintln(out, "YES")
   // output edges
   for i := 0; i < ec; i += 2 {
      u := e[i|1].to
      v := e[i].to
      fmt.Fprint(out, u, " ", v, " ")
      if cho[i] {
         fmt.Fprintln(out, e[i].w)
      } else if cho[i|1] {
         fmt.Fprintln(out, e[i|1].w)
      } else {
         fmt.Fprintln(out, 1000000000000000000)
      }
   }
}
