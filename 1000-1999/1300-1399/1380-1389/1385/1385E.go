package main

import (
   "bufio"
   "fmt"
   "os"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

func nextInt() int {
   var x int
   var c byte
   var neg bool
   // skip non-digit
   for {
       b, err := reader.ReadByte()
       if err != nil {
           return x
       }
       c = b
       if c == '-' || (c >= '0' && c <= '9') {
           break
       }
   }
   if c == '-' {
       neg = true
       b, _ := reader.ReadByte()
       c = b
   }
   for c >= '0' && c <= '9' {
       x = x*10 + int(c-'0')
       b, err := reader.ReadByte()
       if err != nil {
           break
       }
       c = b
   }
   if neg {
       return -x
   }
   return x
}

func main() {
   defer writer.Flush()
   t := nextInt()
   for tc := 0; tc < t; tc++ {
       n := nextInt()
       m := nextInt()
       // directed edges adjacency and indegree
       adj := make([][]int, n+1)
       deg := make([]int, n+1)
       type Edge struct{ u, v, tp int }
       edges := make([]Edge, m)
       for i := 0; i < m; i++ {
           tp := nextInt()
           u := nextInt()
           v := nextInt()
           edges[i] = Edge{u, v, tp}
           if tp == 1 {
               adj[u] = append(adj[u], v)
               deg[v]++
           }
       }
       // topo sort
       ord := make([]int, n+1)
       queue := make([]int, 0, n)
       for i := 1; i <= n; i++ {
           if deg[i] == 0 {
               queue = append(queue, i)
           }
       }
       qi := 0
       cnt := 0
       for qi < len(queue) {
           u := queue[qi]; qi++
           cnt++
           ord[u] = cnt
           for _, v := range adj[u] {
               deg[v]--
               if deg[v] == 0 {
                   queue = append(queue, v)
               }
           }
       }
       if cnt != n {
           fmt.Fprintln(writer, "NO")
       } else {
           fmt.Fprintln(writer, "YES")
           for _, e := range edges {
               if e.tp == 1 {
                   fmt.Fprintf(writer, "%d %d\n", e.u, e.v)
               } else {
                   // orient according to ord
                   if ord[e.u] < ord[e.v] {
                       fmt.Fprintf(writer, "%d %d\n", e.u, e.v)
                   } else {
                       fmt.Fprintf(writer, "%d %d\n", e.v, e.u)
                   }
               }
           }
       }
   }
}
