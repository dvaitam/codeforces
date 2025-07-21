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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   pts := make([][2]float64, n)
   for i := 0; i < n; i++ {
       var xi, yi float64
       fmt.Fscan(reader, &xi, &yi)
       pts[i][0] = xi
       pts[i][1] = yi
   }
   x1 := pts[0][0]
   y0 := pts[0][1]
   x2 := pts[1][0]
   // initial interval on AB
   var L, R float64
   if x1 <= x2 {
       L, R = x1, x2
   } else {
       L, R = x2, x1
   }
   const eps = 1e-9
   for i := 0; i < n; i++ {
       xi := pts[i][0]
       yi := pts[i][1]
       xj := pts[(i+1)%n][0]
       yj := pts[(i+1)%n][1]
       dx := xj - xi
       dy := yj - yi
       A := dx
       B := -dy
       // constraint: A*(y0-yi) + B*(x - xi) <= 0
       if math.Abs(B) < eps {
           // horizontal edge: no x constraint, but must satisfy interior side
           if A*(y0-yi) > 0 {
               fmt.Fprintln(writer, 0)
               return
           }
           continue
       }
       // solve B*(x - xi) <= -A*(y0-yi)
       num := -A * (y0 - yi)
       bound := xi + num/B
       if B > 0 {
           // x <= bound
           if bound < R {
               R = bound
           }
       } else {
           // B < 0: x >= bound
           if bound > L {
               L = bound
           }
       }
       if L > R {
           fmt.Fprintln(writer, 0)
           return
       }
   }
   // count integer x in [L, R]
   start := math.Ceil(L - eps)
   end := math.Floor(R + eps)
   cnt := int(end - start + 1)
   if cnt < 0 {
       cnt = 0
   }
   fmt.Fprintln(writer, cnt)
}
