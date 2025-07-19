package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var T int
   if _, err := fmt.Fscan(in, &T); err != nil {
       return
   }
   for t := 0; t < T; t++ {
       var x [4]float64
       var y [4]float64
       for i := 1; i <= 3; i++ {
           fmt.Fscan(in, &x[i], &y[i])
       }
       ans := 0.0
       if y[1] == y[2] && y[1] > y[3] {
           ans = math.Hypot(x[1]-x[2], y[1]-y[2])
       } else if y[2] == y[3] && y[2] > y[1] {
           ans = math.Hypot(x[2]-x[3], y[2]-y[3])
       } else if y[1] == y[3] && y[3] > y[2] {
           ans = math.Hypot(x[1]-x[3], y[1]-y[3])
       }
       fmt.Fprintf(out, "%.10f\n", ans)
   }
}
