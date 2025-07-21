package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var y1, y2, yw, xb, yb, r float64
   if _, err := fmt.Fscan(reader, &y1, &y2, &yw, &xb, &yb, &r); err != nil {
       return
   }
   // Reflect goal segment across horizontal line y = yw
   yLow := 2*yw - y2
   yHigh := 2*yw - y1
   // Choose a point strictly between to avoid hitting posts exactly
   yPrime := (yLow + yHigh) / 2.0
   // Compute intersection X coordinate on the wall y = yw
   // Line from ball (xb,yb) to reflected point (0,yPrime)
   // xw = xb * (yPrime - yw) / (yPrime - yb)
   denom := yPrime - yb
   if denom == 0 {
       fmt.Println(-1)
       return
   }
   xw := xb * (yPrime - yw) / denom
   if xw <= 0 {
       fmt.Println(-1)
       return
   }
   // Output result with high precision
   fmt.Printf("%.10f\n", xw)
}
