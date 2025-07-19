package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Edge struct {
   z, x, y int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   special := make([]int, n+1)
   for i := 0; i < k; i++ {
       var v int
       fmt.Fscan(reader, &v)
       if v >= 1 && v <= n {
           special[v] = 1
       }
   }
   parent := make([]int, n+1)
   cnt := make([]int, n+1)
   for i := 1; i <= n; i++ {
       parent[i] = i
       cnt[i] = special[i]
   }
   edges := make([]Edge, m)
   for i := 0; i < m; i++ {
       var x, y, z int
       fmt.Fscan(reader, &x, &y, &z)
       edges[i] = Edge{z: z, x: x, y: y}
   }
   sort.Slice(edges, func(i, j int) bool {
       return edges[i].z < edges[j].z
   })

   var find func(int) int
   find = func(u int) int {
       if parent[u] != u {
           parent[u] = find(parent[u])
       }
       return parent[u]
   }
   for _, e := range edges {
       u := find(e.x)
       v := find(e.y)
       if u != v {
           // merge u into v
           parent[u] = v
           cnt[v] += cnt[u]
           if cnt[v] == k {
               // output result
               for i := 0; i < k; i++ {
                   if i > 0 {
                       writer.WriteByte(' ')
                   }
                   fmt.Fprint(writer, e.z)
               }
               writer.WriteByte('\n')
               return
           }
       }
   }
}
