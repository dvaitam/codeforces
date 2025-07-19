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

   var n int
   var p int64
   fmt.Fscan(reader, &n, &p)
   type Point struct {
       dist2 int64
       pop   int64
   }
   pts := make([]Point, n)
   var x, y, z int64
   var maxDist2 int64
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x, &y, &z)
       d2 := x*x + y*y
       pts[i] = Point{dist2: d2, pop: z}
       if d2 > maxDist2 {
           maxDist2 = d2
       }
   }
   needed := int64(1000000) - p
   if needed <= 0 {
       fmt.Fprintln(writer, "0.000000")
       return
   }
   // binary search radius
   st := 0.0
   ed := math.Sqrt(float64(maxDist2))
   for iter := 0; iter < 100; iter++ {
       mid := (st + ed) / 2
       mid2 := mid * mid
       var sum int64
       for _, pt := range pts {
           if float64(pt.dist2) <= mid2 {
               sum += pt.pop
               if sum >= needed {
                   break
               }
           }
       }
       if sum >= needed {
           ed = mid
       } else {
           st = mid
       }
   }
   // final verification
   var total int64
   ed2 := ed * ed
   for _, pt := range pts {
       if float64(pt.dist2) <= ed2 {
           total += pt.pop
       }
   }
   if total < needed {
       fmt.Fprintln(writer, "-1")
   } else {
       fmt.Fprintf(writer, "%.6f\n", ed)
   }
}
