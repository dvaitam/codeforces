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

   v := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &v[i])
   }

   ans := 0
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       x--
       y--
       if v[x] < v[y] {
           ans += v[x]
       } else {
           ans += v[y]
       }
   }

   fmt.Fprintln(out, ans)
}
