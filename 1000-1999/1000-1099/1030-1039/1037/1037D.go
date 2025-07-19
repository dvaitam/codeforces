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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   if n == 0 {
       fmt.Fprintln(writer, "No")
       return
   }
   // root must be 1
   if a[0] != 1 {
       fmt.Fprintln(writer, "No")
       return
   }
   // BFS to compute depth and parent
   depth := make([]int, n+1)
   parent := make([]int, n+1)
   vis := make([]bool, n+1)
   queue := make([]int, 0, n)
   vis[a[0]] = true
   parent[a[0]] = 0
   depth[a[0]] = 0
   queue = append(queue, a[0])
   for i := 0; i < len(queue); i++ {
       u := queue[i]
       for _, v := range adj[u] {
           if !vis[v] {
               vis[v] = true
               parent[v] = u
               depth[v] = depth[u] + 1
               queue = append(queue, v)
           }
       }
   }
   // mark positions
   mark := make([]int, n+1)
   ok := true
   // process sequence in depth groups
   for l := 0; l < n; {
       d := depth[a[l]]
       r := l
       for r+1 < n && depth[a[r+1]] == d {
           r++
       }
       if r+1 < n && depth[a[r+1]] < d {
           ok = false
           break
       }
       for i := l; i <= r; i++ {
           mark[a[i]] = i
           if i > l {
               if mark[parent[a[i]]] < mark[parent[a[i-1]]] {
                   ok = false
                   break
               }
           }
       }
       if !ok {
           break
       }
       l = r + 1
   }
   if ok {
       fmt.Fprintln(writer, "Yes")
   } else {
       fmt.Fprintln(writer, "No")
   }
}
