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

   var n, m int
   fmt.Fscan(reader, &n, &m)

   sol := make([][]int, n+1)
   cnt := 0
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       cnt++
       sol[x] = append(sol[x], cnt)
       sol[y] = append(sol[y], cnt)
   }
   for i := 1; i <= n; i++ {
       if len(sol[i]) == 0 {
           cnt++
           sol[i] = append(sol[i], cnt)
       }
   }
   for i := 1; i <= n; i++ {
       fmt.Fprintln(writer, len(sol[i]))
       for _, id := range sol[i] {
           fmt.Fprintln(writer, i, id)
       }
   }
}
