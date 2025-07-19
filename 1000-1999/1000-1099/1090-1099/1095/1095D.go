package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([][2]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i][0], &a[i][1])
   }

   vis := make([]bool, n+1)
   ans := make([]int, 0, n)
   ans = append(ans, 1)
   vis[1] = true
   for i := 0; i < n; i += 2 {
       u := a[ans[i]][0]
       v := a[ans[i]][1]
       // Determine order to preserve adjacency
       if a[u][0] == v || a[u][1] == v {
           if !vis[u] {
               ans = append(ans, u)
           }
           if !vis[v] {
               ans = append(ans, v)
           }
       } else {
           if !vis[v] {
               ans = append(ans, v)
           }
           if !vis[u] {
               ans = append(ans, u)
           }
       }
       vis[u], vis[v] = true, true
   }

   for i, x := range ans {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, x)
   }
   out.WriteByte('\n')
}
