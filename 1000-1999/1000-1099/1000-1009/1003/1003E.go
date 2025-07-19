package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var N, d, k int
   if _, err := fmt.Fscan(reader, &N, &d, &k); err != nil {
       return
   }
   if d >= N || (k == 1 && N > 2) {
       fmt.Fprintln(writer, "NO")
       return
   }
   if N == 1 && d == 1 {
       fmt.Fprintln(writer, "YES")
       return
   }
   count := make([]int, N+2)
   depth := make([]int, N+2)
   edges := make([][2]int, 0, N-1)

   // Build initial diameter path
   for i := 1; i <= d; i++ {
       u := i
       v := i + 1
       edges = append(edges, [2]int{u, v})
       count[u]++
       count[v]++
       // remaining distance to extend branches
       depth[u] = min(u-1, d-u+1)
       depth[v] = min(u, d-u)
   }
   // Attach remaining nodes
   i := d + 2
   j := 2
   for i <= N {
       for j < i && (count[j] == k || depth[j] == 0) {
           j++
       }
       if i == j {
           fmt.Fprintln(writer, "NO")
           return
       }
       edges = append(edges, [2]int{i, j})
       count[i]++
       count[j]++
       depth[i] = depth[j] - 1
       i++
   }
   if len(edges) != N-1 {
       fmt.Fprintln(writer, "NO")
       return
   }
   fmt.Fprintln(writer, "YES")
   for _, e := range edges {
       fmt.Fprintln(writer, e[0], e[1])
   }
}
