package main

import (
   "bufio"
   "fmt"
   "os"
)

type compEdge struct {
   x, y int
   w    int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k, m int
   if _, err := fmt.Fscan(reader, &n, &k, &m); err != nil {
       return
   }
   // DSU for initial spanning tree
   parent1 := make([]int, n+1)
   for i := 1; i <= n; i++ {
       parent1[i] = i
   }
   var find1 func(int) int
   find1 = func(x int) int {
       if parent1[x] != x {
           parent1[x] = find1(parent1[x])
       }
       return parent1[x]
   }
   union1 := func(x, y int) {
       rx := find1(x)
       ry := find1(y)
       if rx != ry {
           parent1[rx] = ry
       }
   }
   // Graph: adjacency list of tree edges
   type edge2 struct{ to, w int }
   graph := make([][]edge2, n+1)
   // read our k edges (weight=1)
   for i := 0; i < k; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       union1(u, v)
       graph[u] = append(graph[u], edge2{v, 1})
       graph[v] = append(graph[v], edge2{u, 1})
   }
   // read competitor edges
   comp := make([]compEdge, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &comp[i].x, &comp[i].y, &comp[i].w)
   }
   // add competitor edges that connect different components (weight=0)
   for i := 0; i < m; i++ {
       u, v := comp[i].x, comp[i].y
       if find1(u) != find1(v) {
           union1(u, v)
           graph[u] = append(graph[u], edge2{v, 0})
           graph[v] = append(graph[v], edge2{u, 0})
       }
   }
   // BFS to set parent, depth, and initial edge weight u
   parent := make([]int, n+1)
   depth := make([]int, n+1)
   uweight := make([]int, n+1)
   queue := make([]int, 0, n)
   queue = append(queue, 1)
   parent[1] = 0
   depth[1] = 0
   uweight[1] = 0
   for qi := 0; qi < len(queue); qi++ {
       x := queue[qi]
       for _, e := range graph[x] {
           if e.to == parent[x] {
               continue
           }
           parent[e.to] = x
           depth[e.to] = depth[x] + 1
           uweight[e.to] = e.w
           queue = append(queue, e.to)
       }
   }
   // dynamic DSU for path compression
   parent2 := make([]int, n+1)
   for i := 1; i <= n; i++ {
       parent2[i] = i
   }
   var find2 func(int) int
   find2 = func(x int) int {
       if parent2[x] != x {
           parent2[x] = find2(parent2[x])
       }
       return parent2[x]
   }
   ans := make([]int, n+1)
   // process all competitor edges in increasing order
   for i := 0; i < m; i++ {
       x := comp[i].x
       y := comp[i].y
       w := comp[i].w
       x = find2(x)
       y = find2(y)
       for x != y {
           if depth[x] < depth[y] {
               x, y = y, x
           }
           ans[x] = w
           parent2[x] = parent[x]
           x = find2(x)
       }
   }
   // compute answer
   var total int64
   for i := 2; i <= n; i++ {
       if uweight[i] == 1 {
           if ans[i] == 0 {
               fmt.Fprintln(writer, -1)
               return
           }
           total += int64(ans[i])
       }
   }
   fmt.Fprintln(writer, total)
}
