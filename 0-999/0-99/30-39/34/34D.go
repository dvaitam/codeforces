package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, r1, r2 int
   if _, err := fmt.Fscan(reader, &n, &r1, &r2); err != nil {
       return
   }
   // Read old parent representation (skipping old root r1)
   parent := make([]int, n+1)
   for i := 1; i <= n; i++ {
       if i == r1 {
           continue
       }
       fmt.Fscan(reader, &parent[i])
   }
   // Build undirected tree
   adj := make([][]int, n+1)
   for i := 1; i <= n; i++ {
       if i == r1 {
           continue
       }
       p := parent[i]
       adj[i] = append(adj[i], p)
       adj[p] = append(adj[p], i)
   }
   // BFS from new root r2 to compute new parents
   resParent := make([]int, n+1)
   visited := make([]bool, n+1)
   queue := make([]int, 0, n)
   queue = append(queue, r2)
   visited[r2] = true
   for head := 0; head < len(queue); head++ {
       u := queue[head]
       for _, v := range adj[u] {
           if !visited[v] {
               visited[v] = true
               resParent[v] = u
               queue = append(queue, v)
           }
       }
   }
   // Output new parent representation (skipping new root r2)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   first := true
   for i := 1; i <= n; i++ {
       if i == r2 {
           continue
       }
       if !first {
           writer.WriteByte(' ')
       }
       first = false
       fmt.Fprint(writer, resParent[i])
   }
   writer.WriteByte('\n')
}
