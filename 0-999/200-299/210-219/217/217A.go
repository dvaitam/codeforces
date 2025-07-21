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
   fmt.Fscan(reader, &n)
   xs := make([]int, n)
   ys := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &xs[i], &ys[i])
   }

   visited := make([]bool, n)
   var dfs func(int)
   dfs = func(u int) {
       visited[u] = true
       for v := 0; v < n; v++ {
           if !visited[v] && (xs[u] == xs[v] || ys[u] == ys[v]) {
               dfs(v)
           }
       }
   }

   components := 0
   for i := 0; i < n; i++ {
       if !visited[i] {
           components++
           dfs(i)
       }
   }
   // To connect components, need components-1 additional drifts
   fmt.Fprintln(writer, components-1)
}
