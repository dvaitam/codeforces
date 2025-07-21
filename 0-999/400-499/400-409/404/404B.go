package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var a, d float64
   var n int
   fmt.Fscan(reader, &a, &d)
   fmt.Fscan(reader, &n)

   per := 4 * a
   dist := 0.0
   for i := 0; i < n; i++ {
       dist += d
       dist = math.Mod(dist, per)

       var x, y float64
       if dist <= a {
           x = dist
           y = 0
       } else if dist <= 2*a {
           x = a
           y = dist - a
       } else if dist <= 3*a {
           x = a - (dist - 2*a)
           y = a
       } else {
           x = 0
           y = a - (dist - 3*a)
       }

       fmt.Fprintf(writer, "%.6f %.6f\n", x, y)
   }
}
