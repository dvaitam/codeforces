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

   deg := make([]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       deg[u]++
       deg[v]++
   }

   countDeg1, countDeg2, countDegn1 := 0, 0, 0
   for i := 1; i <= n; i++ {
       if deg[i] == 1 {
           countDeg1++
       } else if deg[i] == 2 {
           countDeg2++
       } else if deg[i] == n-1 {
           countDegn1++
       }
   }

   switch {
   case m == n-1 && countDeg1 == 2 && countDeg2 == n-2:
       fmt.Fprintln(writer, "bus topology")
   case m == n && countDeg2 == n:
       fmt.Fprintln(writer, "ring topology")
   case m == n-1 && countDegn1 == 1 && countDeg1 == n-1:
       fmt.Fprintln(writer, "star topology")
   default:
       fmt.Fprintln(writer, "unknown topology")
   }
}
