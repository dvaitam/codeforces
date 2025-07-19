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
   X := make([]int, n)
   Y := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &X[i], &Y[i])
   }
   minX, minY := X[0], Y[0]
   maxX, maxY := X[0], Y[0]
   for i := 1; i < n; i++ {
       if X[i] < minX {
           minX = X[i]
       }
       if X[i] > maxX {
           maxX = X[i]
       }
       if Y[i] < minY {
           minY = Y[i]
       }
       if Y[i] > maxY {
           maxY = Y[i]
       }
   }
   res := 0
   for i := 0; i < n; i++ {
       dx := X[i] - minX
       if maxX - X[i] > dx {
           dx = maxX - X[i]
       }
       dy := Y[i] - minY
       if maxY - Y[i] > dy {
           dy = maxY - Y[i]
       }
       tmp := 2 * (dx + dy)
       if tmp > res {
           res = tmp
       }
   }
   fmt.Fprint(out, res)
   full := 2*((maxX-minX)+(maxY-minY))
   for k := 4; k <= n; k++ {
       fmt.Fprint(out, " ", full)
   }
   fmt.Fprintln(out)
}
