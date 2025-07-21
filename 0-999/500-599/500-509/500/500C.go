package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   w := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &w[i])
   }
   last := make([]int, n+1)
   var ans int64
   for t := 1; t <= m; t++ {
       var x int
       fmt.Fscan(in, &x)
       for i := 1; i <= n; i++ {
           if i != x && last[i] > last[x] {
               ans += w[i]
           }
       }
       last[x] = t
   }
   fmt.Println(ans)
}
