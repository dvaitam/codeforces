package main

import (
   "bufio"
   "container/list"
   "fmt"
   "os"
)

const INF = 1000000000

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, a, b int
   if _, err := fmt.Fscan(in, &n, &m, &a, &b); err != nil {
       return
   }
   a--
   b--
   // read graph
   g := make([][]int, n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--
       v--
       g[u] = append(g[u], v)
   }
   // compute all-pairs shortest dist by BFS
   dist := make([][]int, n)
   for i := 0; i < n; i++ {
       d := make([]int, n)
       for j := 0; j < n; j++ {
           d[j] = INF
       }
       d[i] = 0
       q := list.New()
       q.PushBack(i)
       for q.Len() > 0 {
           u := q.Remove(q.Front()).(int)
           for _, v := range g[u] {
               if d[v] > d[u]+1 {
                   d[v] = d[u] + 1
                   q.PushBack(v)
               }
           }
       }
       dist[i] = d
   }
   var k int
   fmt.Fscan(in, &k)
   // build bus graph: edge v->ti if v on some shortest si->ti path
   busAdj := make([][]int, n)
   for i := 0; i < k; i++ {
       var si, ti int
       fmt.Fscan(in, &si, &ti)
       si--
       ti--
       dsti := dist[si][ti]
       if dsti >= INF {
           continue
       }
       for v := 0; v < n; v++ {
           if dist[si][v] + dist[v][ti] == dsti {
               // can board at v and at worst go to ti
               busAdj[v] = append(busAdj[v], ti)
           }
       }
   }
   // BFS on bus graph from a to b
   dbus := make([]int, n)
   for i := 0; i < n; i++ {
       dbus[i] = INF
   }
   dq := list.New()
   dbus[a] = 0
   dq.PushBack(a)
   for dq.Len() > 0 {
       u := dq.Remove(dq.Front()).(int)
       for _, v := range busAdj[u] {
           if dbus[v] > dbus[u] + 1 {
               dbus[v] = dbus[u] + 1
               dq.PushBack(v)
           }
       }
   }
   if dbus[b] >= INF {
       fmt.Println(-1)
   } else {
       fmt.Println(dbus[b])
   }
}
