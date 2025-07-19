package main

import (
   "bufio"
   "fmt"
   "os"
)

type Edge struct {
   to, id int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, s int
   fmt.Fscan(reader, &n, &m, &s)
   adj := make([][]Edge, n+1)
   iss := make([]bool, m+1)
   for i := 1; i <= m; i++ {
       var t, u, v int
       fmt.Fscan(reader, &t, &u, &v)
       if t == 1 {
           // directed edge
           adj[u] = append(adj[u], Edge{to: v, id: 0})
       } else {
           // undirected, mark index
           iss[i] = true
           adj[u] = append(adj[u], Edge{to: v, id: i})
           adj[v] = append(adj[v], Edge{to: u, id: -i})
       }
   }
   // vis for initial reach via directed edges
   vis := make([]bool, n+1)
   // DFS stack
   stack := make([]int, 0, n)
   vis[s] = true
   stack = append(stack, s)
   for len(stack) > 0 {
       u := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       for _, e := range adj[u] {
           if e.id == 0 && !vis[e.to] {
               vis[e.to] = true
               stack = append(stack, e.to)
           }
       }
   }
   // compute minimal plan ans2, and initial queue for BFS
   ans1 := make([]bool, m+1)
   ans2 := make([]bool, m+1)
   queue := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if vis[i] {
           queue = append(queue, i)
       }
   }
   A2 := len(queue)
   // minimal orientations: orient towards visited
   for _, u := range queue {
       for _, e := range adj[u] {
           if e.id != 0 && !vis[e.to] {
               if e.id > 0 {
                   ans2[e.id] = false // orient to avoid new
               } else {
                   ans2[-e.id] = true
               }
           }
       }
   }
   // maximal plan: BFS expanding both directed and undirected
   A1 := A2
   // use same vis to mark during BFS
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, e := range adj[u] {
           if !vis[e.to] {
               vis[e.to] = true
               if e.id != 0 {
                   if e.id > 0 {
                       ans1[e.id] = true
                   } else {
                       ans1[-e.id] = false
                   }
               }
               queue = append(queue, e.to)
               A1++
           }
       }
   }
   // output max plan
   fmt.Fprintln(writer, A1)
   for i := 1; i <= m; i++ {
       if iss[i] {
           if ans1[i] {
               writer.WriteByte('+')
           } else {
               writer.WriteByte('-')
           }
       }
   }
   writer.WriteByte('\n')
   // output min plan
   fmt.Fprintln(writer, A2)
   for i := 1; i <= m; i++ {
       if iss[i] {
           if ans2[i] {
               writer.WriteByte('+')
           } else {
               writer.WriteByte('-')
           }
       }
   }
   writer.WriteByte('\n')
}
