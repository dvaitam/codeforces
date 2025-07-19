package main

import (
   "bufio"
   "fmt"
   "os"
)

type edge struct {
   id, to int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   cnt := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &cnt[i])
   }
   a0 := make([]int, m+1)
   a1 := make([]int, m+1)
   adj := make([][]edge, n+1)
   for i := 1; i <= m; i++ {
       fmt.Fscan(reader, &a0[i], &a1[i])
       adj[a0[i]] = append(adj[a0[i]], edge{i, a1[i]})
       adj[a1[i]] = append(adj[a1[i]], edge{i, a0[i]})
   }
   siz := make([]int, n+1)
   vis := make([]bool, n+1)
   inEdge := make([]bool, m+1)
   queue := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       siz[i] = len(adj[i])
       if cnt[i] >= siz[i] {
           vis[i] = true
           queue = append(queue, i)
       }
   }
   var stack []int
   qi := 0
   for qi < len(queue) {
       d := queue[qi]
       qi++
       for _, e := range adj[d] {
           if !inEdge[e.id] {
               inEdge[e.id] = true
               stack = append(stack, e.id)
           }
           if vis[e.to] {
               continue
           }
           siz[e.to]--
           if cnt[e.to] >= siz[e.to] {
               vis[e.to] = true
               queue = append(queue, e.to)
           }
       }
   }
   if len(stack) != m {
       fmt.Fprintln(writer, "DEAD")
   } else {
       fmt.Fprintln(writer, "ALIVE")
       for i := m - 1; i >= 0; i-- {
           fmt.Fprint(writer, stack[i])
           if i > 0 {
               fmt.Fprint(writer, " ")
           }
       }
       fmt.Fprintln(writer)
   }
}
