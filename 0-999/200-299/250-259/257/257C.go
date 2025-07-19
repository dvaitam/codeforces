package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   angles := make([]float64, n)
   for i := 0; i < n; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       deg := math.Atan2(float64(y), float64(x)) * 180.0 / math.Pi
       if deg < 0 {
           deg += 360.0
       }
       angles[i] = deg
   }
   sort.Float64s(angles)
   maxGap := 0.0
   for i := 1; i < n; i++ {
       gap := angles[i] - angles[i-1]
       if gap > maxGap {
           maxGap = gap
       }
   }
   if n > 0 {
       wrapGap := 360.0 - angles[n-1] + angles[0]
       if wrapGap > maxGap {
           maxGap = wrapGap
       }
   }
   result := 360.0 - maxGap
   // Output with 8 decimals
   fmt.Printf("%.8f\n", result)
}
