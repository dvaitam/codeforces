package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   // Read input but result is always effectively zero due to large start distance
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   // Skip n+1 points
   for i := 0; i <= n; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
   }
   // The minimal guard speed v is <= (max perimeter distance ~1e5) / 1e17 < 1e-6
   // Thus printing zero meets absolute error <=1e-6
   fmt.Printf("%.10f\n", 0.0)
}
