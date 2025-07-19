package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   G0 := make([][]int, n+1)
   type rel struct{op, u, v int}
   rels := make([]rel, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &rels[i].op, &rels[i].u, &rels[i].v)
       u, v := rels[i].u, rels[i].v
       G0[u] = append(G0[u], v)
       G0[v] = append(G0[v], u)
   }
   col := make([]int, n+1)
   vis := make([]bool, n+1)
   flg := true
   queue := make([]int, n)
   for i := 1; i <= n && flg; i++ {
       if !vis[i] {
           head, tail := 0, 0
           queue[tail] = i; tail++
           vis[i] = true
           col[i] = 1
           for head < tail && flg {
               u := queue[head]; head++
               for _, v := range G0[u] {
                   if vis[v] {
                       if col[v] == col[u] {
                           flg = false
                           break
                       }
                       continue
                   }
                   vis[v] = true
                   col[v] = 3 - col[u]
                   queue[tail] = v; tail++
               }
           }
       }
   }
   if !flg {
       fmt.Fprintln(writer, "NO")
       return
   }
   G := make([][]int, n+1)
   indegree := make([]int, n+1)
   for _, e := range rels {
       u, v := e.u, e.v
       if e.op == 1 {
           if col[u] == 1 {
               G[u] = append(G[u], v)
               indegree[v]++
           } else {
               G[v] = append(G[v], u)
               indegree[u]++
           }
       } else {
           if col[u] == 2 {
               G[u] = append(G[u], v)
               indegree[v]++
           } else {
               G[v] = append(G[v], u)
               indegree[u]++
           }
       }
   }
   topo := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if indegree[i] == 0 {
           topo = append(topo, i)
       }
   }
   x := make([]int, n+1)
   for i := 0; i < len(topo); i++ {
       u := topo[i]
       x[u] = i + 1
       for _, v := range G[u] {
           indegree[v]--
           if indegree[v] == 0 {
               topo = append(topo, v)
           }
       }
   }
   if len(topo) != n {
       fmt.Fprintln(writer, "NO")
       return
   }
   fmt.Fprintln(writer, "YES")
   for i := 1; i <= n; i++ {
       ori := 'L'
       if col[i] == 2 {
           ori = 'R'
       }
       fmt.Fprintf(writer, "%c %d\n", ori, x[i])
   }
}
