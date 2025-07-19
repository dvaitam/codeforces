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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var n int
       fmt.Fscan(reader, &n)
       // n is odd
       half := n / 2
       consider := half + (half % 2)
       // central angle for polygon vertex spacing
       polyAngle := math.Pi * float64(2*n-2) / float64(2*n)
       // circumscribed circle radius
       cRadius := 1.0 / (2 * math.Cos(polyAngle/2))
       // angle to the vertex farthest in square projection
       considerAngle := 2 * math.Pi * float64(consider) / float64(2*n)
       part1 := cRadius * math.Cos(considerAngle/2)
       part2 := cRadius * math.Sin(considerAngle/2)
       ans := (part1 + part2) * math.Sqrt(2)
       fmt.Fprintf(writer, "%.10f\n", ans)
   }
}
