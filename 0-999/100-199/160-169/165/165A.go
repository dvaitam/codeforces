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
   xs := make([]int, n)
   ys := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &xs[i], &ys[i])
   }
   count := 0
   for i := 0; i < n; i++ {
       hasLeft, hasRight, hasUpper, hasLower := false, false, false, false
       xi, yi := xs[i], ys[i]
       for j := 0; j < n; j++ {
           if j == i {
               continue
           }
           xj, yj := xs[j], ys[j]
           if xj > xi && yj == yi {
               hasRight = true
           }
           if xj < xi && yj == yi {
               hasLeft = true
           }
           if xj == xi && yj > yi {
               hasUpper = true
           }
           if xj == xi && yj < yi {
               hasLower = true
           }
       }
       if hasLeft && hasRight && hasUpper && hasLower {
           count++
       }
   }
   fmt.Fprintln(writer, count)
}
