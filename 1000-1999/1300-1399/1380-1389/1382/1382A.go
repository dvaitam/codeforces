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

   var t int
   fmt.Fscan(in, &t)
   for t > 0 {
       t--
       var n, m int
       fmt.Fscan(in, &n, &m)
       seen := make(map[int]bool, n)
       for i := 0; i < n; i++ {
           var v int
           fmt.Fscan(in, &v)
           seen[v] = true
       }
       res := -1
       for i := 0; i < m; i++ {
           var v int
           fmt.Fscan(in, &v)
           if res == -1 && seen[v] {
               res = v
           }
       }
       if res != -1 {
           fmt.Fprintln(out, "YES")
           fmt.Fprintln(out, 1, res)
       } else {
           fmt.Fprintln(out, "NO")
       }
   }
}
