package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 1000000000

// Edge for flow graph
type Edge struct {
   to, rev, cap int
}

func min(a, b int) int { if a < b { return a }; return b }

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   eee := make([][]int, n)
   for i := 0; i < n; i++ {
       var k int
       fmt.Fscan(in, &k)
       eee[i] = make([]int, k)
       for j := 0; j < k; j++ {
           fmt.Fscan(in, &eee[i][j])
           eee[i][j]--
       }
   }
   cs := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &cs[i])
   }
   // Bipartite matching p: left->right, q: right->left
   p := make([]int, n)
   q := make([]int, n)
   for i := 0; i < n; i++ { p[i], q[i] = -1, -1 }
   was := make([]bool, n)
   var dfsc func(int) bool
   dfsc = func(v int) bool {
       was[v] = true
       for _, u := range eee[v] {
           if q[u] == -1 {
               q[u] = v
               p[v] = u
               return true
           }
       }
       for _, u := range eee[v] {
           w := q[u]
           if !was[w] {
               if dfsc(w) {
                   q[u] = v
                   p[v] = u
                   return true
               }
           }
       }
       return false
   }
   now := 0
   for now < n {
       for i := range was { was[i] = false }
       for i := 0; i < n; i++ {
           if p[i] == -1 && !was[i] {
               if dfsc(i) {
                   now++
               }
           }
       }
   }
   // build flow graph
   N := n + 2
   S, T := n, n+1
   G := make([][]Edge, N)
   addEdge := func(u, v, c int) {
       G[u] = append(G[u], Edge{v, len(G[v]), c})
       G[v] = append(G[v], Edge{u, len(G[u]) - 1, 0})
   }
   sumNeg := 0
   for i := 0; i < n; i++ {
       if cs[i] >= 0 {
           addEdge(i, T, cs[i])
       } else {
           sumNeg += cs[i]
           addEdge(S, i, -cs[i])
       }
   }
   for i := 0; i < n; i++ {
       for _, j := range eee[i] {
           if q[j] != i {
               addEdge(i, q[j], INF)
           }
       }
   }
   // Dinic with scaling
   level := make([]int, N)
   iter := make([]int, N)
   var bfs = func(mn int) bool {
       for i := range level { level[i] = -1 }
       queue := make([]int, 0, N)
       level[S] = 0
       queue = append(queue, S)
       for qi := 0; qi < len(queue); qi++ {
           v := queue[qi]
           for _, e := range G[v] {
               if e.cap >= mn && level[e.to] < 0 {
                   level[e.to] = level[v] + 1
                   queue = append(queue, e.to)
               }
           }
       }
       return level[T] >= 0
   }
   var dfsf func(v, upTo, mn int) int
   dfsf = func(v, upTo, mn int) int {
       if v == T {
           return upTo
       }
       res := 0
       for ; upTo >= mn && iter[v] < len(G[v]); iter[v]++ {
           e := &G[v][iter[v]]
           if e.cap < mn || level[e.to] != level[v]+1 {
               continue
           }
           f := dfsf(e.to, min(upTo, e.cap), mn)
           if f > 0 {
               e.cap -= f
               G[e.to][e.rev].cap += f
               upTo -= f
               res += f
           }
       }
       return res
   }
   flow := 0
   for mn := 1 << 20; mn > 0; mn >>= 1 {
       for {
           if !bfs(mn) {
               break
           }
           for i := range iter { iter[i] = 0 }
           f := dfsf(S, INF, mn)
           if f == 0 {
               break
           }
           flow += f
       }
   }
   fmt.Fprintln(out, flow+sumNeg)
}
