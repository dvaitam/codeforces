package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type edge struct{ to, rev, cap int }

var (
   a       []int
   vSorted []int
   ans     []int
)

func solve(s, e int, vtx []int, edg [][2]int) {
   if s == e {
       for _, x := range vtx {
           ans[x] = vSorted[s]
       }
       return
   }
   m := (s + e) / 2
   N := len(vtx)
   S, T := N, N+1
   // build graph
   g := make([][]edge, N+2)
   addEdge := func(u, v, c int) {
       g[u] = append(g[u], edge{v, len(g[v]), c})
       g[v] = append(g[v], edge{u, len(g[u]) - 1, 0})
   }
   for i, vx := range vtx {
       if a[vx] <= vSorted[m] {
           addEdge(S, i, 1)
       } else {
           addEdge(i, T, 1)
       }
   }
   inf := N + 1
   for _, uv := range edg {
       u, v := uv[0], uv[1]
       // find indices in vtx
       iu := sort.SearchInts(vtx, u)
       iv := sort.SearchInts(vtx, v)
       if iu < N && vtx[iu] == u && iv < N && vtx[iv] == v {
           // add infinite capacity edge from v to u
           addEdge(iv, iu, inf)
       }
   }
   // Dinic maxflow
   level := make([]int, N+2)
   iter := make([]int, N+2)
   bfs := func() bool {
       for i := range level {
           level[i] = -1
       }
       queue := make([]int, 0, N+2)
       level[S] = 0
       queue = append(queue, S)
       for i := 0; i < len(queue); i++ {
           u := queue[i]
           for _, e := range g[u] {
               if e.cap > 0 && level[e.to] < 0 {
                   level[e.to] = level[u] + 1
                   queue = append(queue, e.to)
               }
           }
       }
       return level[T] >= 0
   }
   // Dinic implementation
   var dfs func(u, f int) int
   flow := 0
   maxCap := N + 5
   dfs = func(u, f int) int {
       if u == T {
           return f
       }
       for iter[u] < len(g[u]) {
           e := &g[u][iter[u]]
           if e.cap > 0 && level[u] < level[e.to] {
               d := dfs(e.to, min(f, e.cap))
               if d > 0 {
                   e.cap -= d
                   g[e.to][e.rev].cap += d
                   return d
               }
           }
           iter[u]++
       }
       return 0
   }
   for bfs() {
       for i := range iter {
           iter[i] = 0
       }
       for {
           pushed := dfs(S, maxCap)
           if pushed == 0 {
               break
           }
           flow += pushed
       }
   }
   // reachable from S
   reach := make([]bool, N+2)
   queue := []int{S}
   reach[S] = true
   for i := 0; i < len(queue); i++ {
       u := queue[i]
       for _, e := range g[u] {
           if e.cap > 0 && !reach[e.to] {
               reach[e.to] = true
               queue = append(queue, e.to)
           }
       }
   }
   // partition
   var l, r []int
   var le, re [][2]int
   for i, vx := range vtx {
       if reach[i] {
           l = append(l, vx)
       } else {
           r = append(r, vx)
       }
   }
   for _, uv := range edg {
       u, v := uv[0], uv[1]
       iu := sort.SearchInts(vtx, u)
       iv := sort.SearchInts(vtx, v)
       if iu < N && iv < N {
           on := reach[iu]
           ov := reach[iv]
           if on != ov {
               continue
           }
           if on {
               le = append(le, uv)
           } else {
               re = append(re, uv)
           }
       }
   }
   solve(s, m, l, le)
   solve(m+1, e, r, re)
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m int
   fmt.Fscan(in, &n, &m)
   a = make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   b := make([][2]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &b[i][0], &b[i][1])
       b[i][0]--
       b[i][1]--
   }
   vSorted = make([]int, n)
   copy(vSorted, a)
   sort.Ints(vSorted)
   vSorted = unique(vSorted)
   ans = make([]int, n)
   vtx := make([]int, n)
   for i := 0; i < n; i++ {
       vtx[i] = i
   }
   solve(0, len(vSorted)-1, vtx, b)
   for i, val := range ans {
       if i > 0 {
           fmt.Fprint(out, " ")
       }
       fmt.Fprint(out, val)
   }
}

func unique(a []int) []int {
   j := 0
   for i := 0; i < len(a); i++ {
       if i == 0 || a[i] != a[i-1] {
           a[j] = a[i]
           j++
       }
   }
   return a[:j]
}
