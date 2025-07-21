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
   adj := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   comp := make([]int, n+1)
   color := make([]int, n+1)
   visited := make([]bool, n+1)
   isBip := make([]bool, 0, n)
   compID := 0
   // BFS per component
   for v := 1; v <= n; v++ {
       if visited[v] {
           continue
       }
       // new component
       isBip = append(isBip, true)
       visited[v] = true
       comp[v] = compID
       color[v] = 0
       queue := []int{v}
       for head := 0; head < len(queue); head++ {
           u := queue[head]
           for _, w := range adj[u] {
               if !visited[w] {
                   visited[w] = true
                   comp[w] = compID
                   color[w] = color[u] ^ 1
                   queue = append(queue, w)
               } else if comp[w] == compID && color[w] == color[u] {
                   isBip[compID] = false
               }
           }
       }
       compID++
   }
   var q int
   fmt.Fscan(reader, &q)
   for i := 0; i < q; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       ans := "No"
       if u != v && comp[u] == comp[v] {
           cid := comp[u]
           if !isBip[cid] || color[u] != color[v] {
               ans = "Yes"
           }
       }
       writer.WriteString(ans)
       writer.WriteByte('\n')
   }
}
