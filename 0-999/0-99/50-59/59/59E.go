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
   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   // read graph
   adj := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       adj[x] = append(adj[x], y)
       adj[y] = append(adj[y], x)
   }
   // read forbidden triples
   n1 := int64(n) + 1
   forbidden := make(map[int64]bool, k)
   for i := 0; i < k; i++ {
       var a, b, c int
       fmt.Fscan(reader, &a, &b, &c)
       key := (int64(a)*n1 + int64(b)) * n1 + int64(c)
       forbidden[key] = true
   }
   // BFS on states (prev, curr)
   // state key = prev*n1 + curr
   initial := int64(0)*n1 + 1
   dist := make(map[int64]int, 0)
   parent := make(map[int64]int64, 0)
   queue := make([]int64, 0, 10000)
   dist[initial] = 0
   parent[initial] = -1
   queue = append(queue, initial)
   var ansKey int64 = -1
   for head := 0; head < len(queue); head++ {
       cur := queue[head]
       d := dist[cur]
       u := int(cur / n1)
       v := int(cur % n1)
       if v == n {
           ansKey = cur
           break
       }
       for _, w := range adj[v] {
           if u != 0 {
               tkey := cur*n1 + int64(w)
               if forbidden[tkey] {
                   continue
               }
           }
           next := int64(v)*n1 + int64(w)
           if _, seen := dist[next]; !seen {
               dist[next] = d + 1
               parent[next] = cur
               queue = append(queue, next)
           }
       }
   }
   if ansKey < 0 {
       fmt.Fprintln(writer, -1)
       return
   }
   // reconstruct path
   d := dist[ansKey]
   path := make([]int, 0, d+1)
   for key := ansKey; key >= 0; key = parent[key] {
       v := int(key % n1)
       path = append(path, v)
   }
   // reverse
   for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
       path[i], path[j] = path[j], path[i]
   }
   // output
   fmt.Fprintln(writer, d)
   for i, v := range path {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}
