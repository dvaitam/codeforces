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
   for {
       if _, err := fmt.Fscan(in, &n); err != nil {
           break
       }
       a := make([]int, n+1)
       sum := 0
       for i := 1; i <= n; i++ {
           fmt.Fscan(in, &a[i])
           sum += a[i]
       }
       vis := make([]bool, n+1)
       d := sum / (n/2)
       for i := 1; i <= n; i++ {
           if vis[i] {
               continue
           }
           for j := i + 1; j <= n; j++ {
               if !vis[j] && a[i]+a[j] == d {
                   vis[i] = true
                   vis[j] = true
                   fmt.Fprintln(out, i, j)
                   break
               }
           }
       }
   }
}
