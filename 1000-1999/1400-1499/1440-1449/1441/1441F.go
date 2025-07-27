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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for t > 0 {
       t--
       var n, m int
       fmt.Fscan(reader, &n, &m)
       // skip coordinates
       for i := 0; i < n; i++ {
           var x, y int
           fmt.Fscan(reader, &x, &y)
       }
       // adjacency
       adj := make([][]int, n)
       for i := 0; i < m; i++ {
           var u, v int
           fmt.Fscan(reader, &u, &v)
           u--
           v--
           adj[u] = append(adj[u], v)
           adj[v] = append(adj[v], u)
       }
       // if odd vertices no perfect matching
       if n&1 == 1 {
           fmt.Fprintln(writer, 0)
           continue
       }
       deg := make([]int, n)
       removed := make([]bool, n)
       for i := 0; i < n; i++ {
           deg[i] = len(adj[i])
       }
       // queue of leaves
       q := make([]int, 0, n)
       for i := 0; i < n; i++ {
           if deg[i] == 1 {
               q = append(q, i)
           }
       }
       removedCount := 0
       // process leaves
       for head := 0; head < len(q); head++ {
           v := q[head]
           if removed[v] || deg[v] != 1 {
               continue
           }
           // find neighbor
           var u int = -1
           for _, w := range adj[v] {
               if !removed[w] {
                   u = w
                   break
               }
           }
           if u < 0 {
               continue
           }
           // match v-u and remove
           removed[v] = true
           removed[u] = true
           removedCount += 2
           // update neighbors of u
           for _, w := range adj[u] {
               if !removed[w] {
                   deg[w]--
                   if deg[w] == 1 {
                       q = append(q, w)
                   }
               }
           }
       }
       if removedCount == n {
           fmt.Fprintln(writer, 1)
       } else {
           fmt.Fprintln(writer, 0)
       }
   }
}
