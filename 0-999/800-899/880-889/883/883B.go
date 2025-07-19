package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// PQEntry is an entry in the max-heap by w
type PQEntry struct {
   w, id int
}

// MaxHeap implements heap.Interface for PQEntry max-heap
type MaxHeap []PQEntry

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].w > h[j].w }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) {
   *h = append(*h, x.(PQEntry))
}
func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   if n < k {
       fmt.Fprintln(out, -1)
       return
   }
   a := make([]int, n)
   flag := make([]bool, k+1)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
       if a[i] > 0 && a[i] <= k {
           flag[a[i]] = true
       }
   }
   // graph reversed and original
   graph := make([][]int, n)
   graph1 := make([][]int, n)
   d := make([]int, n)
   d1 := make([]int, n)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       x--
       y--
       // reversed edge y->x
       graph[y] = append(graph[y], x)
       d[x]++
       // original edge x->y
       graph1[x] = append(graph1[x], y)
       d1[y]++
   }
   // dp forward on reversed graph
   dp := make([]int, n)
   D := make([]int, n)
   // initial dp values
   for i := 0; i < n; i++ {
       dp[i] = a[i]
       D[i] = d[i]
   }
   for i := 0; i < n; i++ {
       if dp[i] == 0 {
           dp[i] = 1
       }
   }
   // topo order
   q := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if d[i] == 0 {
           q = append(q, i)
       }
   }
   for idx := 0; idx < len(q); idx++ {
       u := q[idx]
       for _, v := range graph[u] {
           if dp[v] < dp[u]+1 {
               dp[v] = dp[u] + 1
           }
           d[v]--
           if d[v] == 0 {
               q = append(q, v)
           }
       }
   }
   if len(q) != n {
       fmt.Fprintln(out, -1)
       return
   }
   for i := 0; i < n; i++ {
       if a[i] > 0 && dp[i] != a[i] {
           fmt.Fprintln(out, -1)
           return
       }
       if dp[i] > k {
           fmt.Fprintln(out, -1)
           return
       }
   }
   // dp1 backward on original graph
   dp1 := make([]int, n)
   // reset queue
   q = q[:0]
   for i := 0; i < n; i++ {
       if a[i] == 0 {
           dp1[i] = k
       } else {
           dp1[i] = a[i]
       }
       if d1[i] == 0 {
           q = append(q, i)
       }
   }
   for idx := 0; idx < len(q); idx++ {
       u := q[idx]
       for _, v := range graph1[u] {
           if dp1[v] > dp1[u]-1 {
               dp1[v] = dp1[u] - 1
           }
           d1[v]--
           if d1[v] == 0 {
               q = append(q, v)
           }
       }
   }
   // collect free nodes
   type Entry struct{ id, w, w1 int }
   free := make([]Entry, 0, n)
   for i := 0; i < n; i++ {
       if a[i] == 0 {
           free = append(free, Entry{i, dp[i], dp1[i]})
       }
   }
   // sort by w1 desc
   sort.Slice(free, func(i, j int) bool {
       return free[i].w1 > free[j].w1
   })
   // assign values
   mh := &MaxHeap{}
   heap.Init(mh)
   now := 0
   for i := k; i >= 1; i-- {
       if i <= k && (i < len(flag) && !flag[i]) {
           for now < len(free) && free[now].w1 >= i {
               heap.Push(mh, PQEntry{free[now].w, free[now].id})
               now++
           }
           // remove invalid
           for mh.Len() > 0 {
               top := (*mh)[0]
               if top.w > i {
                   heap.Pop(mh)
                   continue
               }
               break
           }
           if mh.Len() == 0 {
               fmt.Fprintln(out, -1)
               return
           }
           e := heap.Pop(mh).(PQEntry)
           a[e.id] = i
       }
   }
   // rebuild dp with assigned a
   for i := 0; i < n; i++ {
       d[i] = D[i]
       dp[i] = a[i]
       if dp[i] == 0 {
           dp[i] = 1
       }
   }
   // final topo
   q = q[:0]
   for i := 0; i < n; i++ {
       if d[i] == 0 {
           q = append(q, i)
       }
   }
   for idx := 0; idx < len(q); idx++ {
       u := q[idx]
       for _, v := range graph[u] {
           if dp[v] < dp[u]+1 {
               dp[v] = dp[u] + 1
           }
           d[v]--
           if d[v] == 0 {
               q = append(q, v)
           }
       }
   }
   // output
   for i := 0; i < n; i++ {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, dp[i])
   }
   out.WriteByte('\n')
}
