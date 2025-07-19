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
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
      fmt.Fscan(reader, &a[i])
   }
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
      var u, v int
      fmt.Fscan(reader, &u, &v)
      adj[u] = append(adj[u], v)
      adj[v] = append(adj[v], u)
   }
   parent := make([]int, n+1)
   queue := make([]int, 0, n)
   queue = append(queue, 1)
   parent[1] = 0
   for i := 0; i < len(queue); i++ {
      u := queue[i]
      for _, v := range adj[u] {
         if v != parent[u] {
            parent[v] = u
            queue = append(queue, v)
         }
      }
   }
   sum1 := make([]int64, n+1)
   for i := len(queue) - 1; i >= 0; i-- {
      u := queue[i]
      s := a[u]
      for _, v := range adj[u] {
         if v != parent[u] && sum1[v] > 0 {
            s += sum1[v]
         }
      }
      sum1[u] = s
   }
   maxx := sum1[1]
   for i := 2; i <= n; i++ {
      if sum1[i] > maxx {
         maxx = sum1[i]
      }
   }
   sum2 := make([]int64, n+1)
   var ans int64
   for i := len(queue) - 1; i >= 0; i-- {
      u := queue[i]
      s := a[u]
      for _, v := range adj[u] {
         if v != parent[u] && sum2[v] > 0 {
            s += sum2[v]
         }
      }
      if s == maxx {
         ans++
         sum2[u] = 0
      } else {
         sum2[u] = s
      }
   }
   fmt.Fprintln(writer, ans*maxx, ans)
}
