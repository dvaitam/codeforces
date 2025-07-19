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
   a := make([]int, n+2)
   b := make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   // f[x] = steps from x to reach n, pre[x] = next position after x
   f := make([]int, n+2)
   pre := make([]int, n+2)
   for i := 0; i <= n; i++ {
       f[i] = -1
       pre[i] = 0
   }
   f[n] = 0
   pre[n] = n + 1
   l, r := n, n
   for r > 0 {
       // find new left boundary: first l where f[l] != f[r]
       for l > 0 && f[l] == f[r] {
           l--
       }
       // process range (l, r]
       for i := r; i > l; i-- {
           base := i + b[i]
           if base < 0 || base > n {
               continue
           }
           low := base - a[base]
           if low < 0 {
               low = 0
           }
           for j := low; j <= base; j++ {
               if pre[j] == 0 {
                   f[j] = f[r] + 1
                   pre[j] = i
               }
           }
       }
       if l < 0 || pre[l] == 0 {
           fmt.Fprintln(writer, -1)
           return
       }
       r = l
   }
   // output result
   fmt.Fprintln(writer, f[0])
   // build path
   var path []int
   var dfs func(int)
   dfs = func(x int) {
       if pre[x] <= n {
           dfs(pre[x])
           path = append(path, x)
       }
   }
   dfs(0)
   for _, x := range path {
       fmt.Fprint(writer, x, " ")
   }
   fmt.Fprintln(writer)
}
