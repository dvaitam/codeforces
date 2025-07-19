package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// Env represents an interval with start s, end t, jump destination d, and weight w.
type Env struct {
   s, t, d int
   w       int64
}

// EnvPQ implements a max-heap on Env based on weight then d.
type EnvPQ []Env

func (h EnvPQ) Len() int { return len(h) }
// Less defines a max-heap: higher w (or, if equal, higher d) has higher priority (is "smaller").
func (h EnvPQ) Less(i, j int) bool {
   if h[i].w != h[j].w {
       return h[i].w > h[j].w
   }
   return h[i].d > h[j].d
}
func (h EnvPQ) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *EnvPQ) Push(x interface{}) {
   *h = append(*h, x.(Env))
}
func (h *EnvPQ) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

func minInt64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   envs := make([]Env, k)
   for i := 0; i < k; i++ {
       var s, t, d int
       var w int64
       fmt.Fscan(in, &s, &t, &d, &w)
       envs[i] = Env{s: s, t: t, d: d, w: w}
   }
   // sort by start
   sort.Slice(envs, func(i, j int) bool {
       return envs[i].s < envs[j].s
   })
   // prepare next position and weight arrays
   md := make([]int, n+2)
   mv := make([]int64, n+2)
   // priority queue of Env
   pq := &EnvPQ{}
   heap.Init(pq)
   head := 0
   for i := 1; i <= n; i++ {
       for head < k && envs[head].s <= i {
           heap.Push(pq, envs[head])
           head++
       }
       for pq.Len() > 0 && (*pq)[0].t < i {
           heap.Pop(pq)
       }
       if pq.Len() == 0 {
           md[i] = i + 1
           mv[i] = 0
       } else {
           top := (*pq)[0]
           md[i] = top.d + 1
           mv[i] = top.w
       }
   }
   // dp arrays for two parities
   dp0 := make([]int64, n+3)
   dp1 := make([]int64, n+3)
   // base case j = 0
   for i := n; i >= 1; i-- {
       dp0[i] = dp0[md[i]] + mv[i]
   }
   // iterate skips j from 1 to m
   for j := 1; j <= m; j++ {
       if j&1 == 1 {
           // fill dp1, read dp0
           for i := n; i >= 1; i-- {
               v1 := dp0[i+1]
               v2 := dp1[md[i]] + mv[i]
               dp1[i] = minInt64(v1, v2)
           }
       } else {
           // fill dp0, read dp1
           for i := n; i >= 1; i-- {
               v1 := dp1[i+1]
               v2 := dp0[md[i]] + mv[i]
               dp0[i] = minInt64(v1, v2)
           }
       }
   }
   // result
   var ans int64
   if m&1 == 1 {
       ans = dp1[1]
   } else {
       ans = dp0[1]
   }
   fmt.Println(ans)
}
