package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const Log = 19

type Edge struct {
   u, v, w, c, id int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, s int
   fmt.Fscan(reader, &n, &m)
   edges := make([]Edge, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &edges[i].w)
   }
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &edges[i].c)
   }
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &edges[i].u, &edges[i].v)
       edges[i].u--
       edges[i].v--
       edges[i].id = i
   }
   fmt.Fscan(reader, &s)

   // sort by weight
   sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })
   to := make([]int, m)
   for i := 0; i < m; i++ {
       to[edges[i].id] = i
   }
   // DSU
   parent := make([]int, n)
   for i := 0; i < n; i++ {
       parent[i] = i
   }
   var find = func(x int) int {
       for parent[x] != x {
           parent[x] = parent[parent[x]]
           x = parent[x]
       }
       return x
   }
   totw := int64(0)
   minc := int(1e9)
   pc := -1
   vis := make([]bool, m)
   // adjacency list for MST
   adj := make([][]struct{to, w int}, n)
   for i := 0; i < m; i++ {
       u := edges[i].u
       v := edges[i].v
       fu := find(u)
       fv := find(v)
       if fu != fv {
           parent[fu] = fv
           totw += int64(edges[i].w)
           if edges[i].c < minc {
               minc = edges[i].c
               pc = i
           }
           vis[i] = true
           adj[u] = append(adj[u], struct{to, w int}{v, edges[i].w})
           adj[v] = append(adj[v], struct{to, w int}{u, edges[i].w})
       }
   }
   // BFS for LCA
   depth := make([]int, n)
   parentUp := make([][]int, Log)
   maxUp := make([][]int, Log)
   for j := 0; j < Log; j++ {
       parentUp[j] = make([]int, n)
       maxUp[j] = make([]int, n)
   }
   // root at 0
   queue := make([]int, 0, n)
   queue = append(queue, 0)
   depth[0] = 0
   parentUp[0][0] = -1
   // BFS
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, e := range adj[u] {
           v := e.to
           if v == parentUp[0][u] {
               continue
           }
           depth[v] = depth[u] + 1
           parentUp[0][v] = u
           maxUp[0][v] = e.w
           queue = append(queue, v)
       }
   }
   // binary lifting
   for j := 1; j < Log; j++ {
       for i := 0; i < n; i++ {
           p := parentUp[j-1][i]
           if p < 0 {
               parentUp[j][i] = -1
               maxUp[j][i] = maxUp[j-1][i]
           } else {
               parentUp[j][i] = parentUp[j-1][p]
               if parentUp[j][i] >= 0 && maxUp[j-1][p] > maxUp[j-1][i] {
                   maxUp[j][i] = maxUp[j-1][p]
               } else {
                   maxUp[j][i] = maxUp[j-1][i]
               }
           }
       }
   }
   // function to get max edge on path
   getMax := func(u, v int) int {
       var res int
       if depth[u] < depth[v] {
           u, v = v, u
       }
       // lift u
       dd := depth[u] - depth[v]
       for j := 0; j < Log; j++ {
           if dd&(1<<j) != 0 {
               if maxUp[j][u] > res {
                   res = maxUp[j][u]
               }
               u = parentUp[j][u]
           }
       }
       if u == v {
           return res
       }
       for j := Log - 1; j >= 0; j-- {
           if parentUp[j][u] != parentUp[j][v] {
               if maxUp[j][u] > res {
                   res = maxUp[j][u]
               }
               if maxUp[j][v] > res {
                   res = maxUp[j][v]
               }
               u = parentUp[j][u]
               v = parentUp[j][v]
           }
       }
       // last step
       if maxUp[0][u] > res {
           res = maxUp[0][u]
       }
       if maxUp[0][v] > res {
           res = maxUp[0][v]
       }
       return res
   }
   // evaluate best
   ans := totw - int64(s/minc)
   pos := -1
   for i := 0; i < m; i++ {
       if !vis[i] && edges[i].c < minc {
           wmx := getMax(edges[i].u, edges[i].v)
           alt := totw - int64(wmx) + int64(edges[i].w) - int64(s/edges[i].c)
           if alt < ans {
               ans = alt
               pos = i
           }
       }
   }
   // output
   fmt.Fprintln(writer, ans)
   if pos < 0 {
       // no replacement
       for orig := 0; orig < m; orig++ {
           idx := to[orig]
           if vis[idx] {
               if idx == pc {
                   fmt.Fprintf(writer, "%d %d\n", orig+1, edges[idx].w - s/minc)
               } else {
                   fmt.Fprintf(writer, "%d %d\n", orig+1, edges[idx].w)
               }
           }
       }
   } else {
       // find edge to remove
       u := edges[pos].u
       v := edges[pos].v
       U, V := 0, 0
       W := 0
       for u != v {
           if depth[u] < depth[v] {
               u, v = v, u
           }
           if maxUp[0][u] > W {
               W = maxUp[0][u]
               U = u
               V = parentUp[0][u]
           }
           u = parentUp[0][u]
       }
       if U < V {
           U, V = V, U
       }
       for orig := 0; orig < m; orig++ {
           idx := to[orig]
           if vis[idx] && !(max(edges[idx].u, edges[idx].v) == U && min(edges[idx].u, edges[idx].v) == V) {
               fmt.Fprintf(writer, "%d %d\n", orig+1, edges[idx].w)
           } else if idx == pos {
               fmt.Fprintf(writer, "%d %d\n", orig+1, edges[idx].w - s/edges[idx].c)
           }
       }
   }
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}
