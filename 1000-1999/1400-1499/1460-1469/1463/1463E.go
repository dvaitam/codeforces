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
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   f := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &f[i])
   }
   nxt := make([]int, n+1)
   vis := make([]bool, n+1)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       nxt[x] = y
       vis[y] = true
   }
   mp := make([][]int, n+1)
   id := make([]int, n+1)
   indeg := make([]int, n+1)
   rt := 0
   // Identify chains and build dependencies
   for i := 1; i <= n; i++ {
       if !vis[i] {
           x := i
           for x != 0 {
               id[x] = i
               if f[x] != 0 && id[f[x]] != i {
                   mp[f[x]] = append(mp[f[x]], i)
                   indeg[i]++
               }
               x = nxt[x]
           }
           if indeg[i] == 0 {
               rt = i
           }
       }
   }
   if rt == 0 {
       fmt.Fprint(out, 0)
       return
   }
   // Topological traversal
   queue := []int{rt}
   ans := make([]int, 0, n)
   for front := 0; front < len(queue); front++ {
       head := queue[front]
       x := head
       for x != 0 {
           ans = append(ans, x)
           for _, ci := range mp[x] {
               indeg[ci]--
               if indeg[ci] == 0 {
                   queue = append(queue, ci)
               }
           }
           x = nxt[x]
       }
   }
   // Output result
   if len(ans) == n {
       for _, v := range ans {
           fmt.Fprint(out, v, " ")
       }
   } else {
       fmt.Fprint(out, 0)
   }
}
