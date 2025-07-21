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
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   if m < 3 || n < m || n > 2*m {
       fmt.Fprintln(writer, -1)
       return
   }
   // Generate m convex hull points: regular polygon on circle
   const R = 50000000.0
   pts := make([][2]int, 0, n)
   for i := 0; i < m; i++ {
       ang := 2 * math.Pi * float64(i) / float64(m)
       x := int(math.Round(R * math.Cos(ang)))
       y := int(math.Round(R * math.Sin(ang)))
       pts = append(pts, [2]int{x, y})
   }
   // Generate interior points on a small parabola
   for j := 1; j <= n-m; j++ {
       x := j
       y := j*j + 1
       pts = append(pts, [2]int{x, y})
   }
   // Output points
   for _, p := range pts {
       fmt.Fprintf(writer, "%d %d\n", p[0], p[1])
   }
}
