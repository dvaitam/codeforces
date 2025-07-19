package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct { v, w int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, q int
   if _, err := fmt.Fscan(reader, &n, &q); err != nil {
       return
   }
   d := make([]int, n)
   edges := make([][2]int, n-1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       d[u] ^= 1
       d[v] ^= 1
       edges[i] = [2]int{u, v}
   }
   adj := make([][]pair, n)
   for i := 0; i < q; i++ {
       var u, v, w int
       fmt.Fscan(reader, &u, &v, &w)
       u--
       v--
       adj[u] = append(adj[u], pair{v, w})
       adj[v] = append(adj[v], pair{u, w})
   }
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       a[i] = -1
   }
   // BFS for each component
   for i := 0; i < n; i++ {
       if a[i] != -1 {
           continue
       }
       queue := make([]int, 0, n)
       head := 0
       queue = append(queue, i)
       a[i] = 0
       for head < len(queue) {
           u := queue[head]
           head++
           b[u] = i
           for _, p := range adj[u] {
               v, w := p.v, p.w
               if a[v] == -1 {
                   a[v] = a[u] ^ w
                   queue = append(queue, v)
               } else if a[v] != (a[u] ^ w) {
                   fmt.Fprintln(writer, "No")
                   return
               }
           }
       }
   }
   ans := 0
   cnt := make([]int, n)
   for i := 0; i < n; i++ {
       if d[i] == 1 {
           ans ^= a[i]
           cnt[b[i]]++
       }
   }
   // adjust one component if needed
   found := false
   for comp := 0; comp < n; comp++ {
       if cnt[comp]%2 == 1 {
           for j := 0; j < n; j++ {
               if b[j] == comp {
                   a[j] ^= ans
               }
           }
           found = true
       }
       if found {
           break
       }
   }
   // output
   fmt.Fprintln(writer, "Yes")
   for i, e := range edges {
       if i > 0 {
           writer.WriteByte(' ')
       }
       u, v := e[0], e[1]
       fmt.Fprint(writer, a[u]^a[v])
   }
   fmt.Fprintln(writer)
}
