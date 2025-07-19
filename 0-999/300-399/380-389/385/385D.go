package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var l, r float64
   if _, err := fmt.Fscan(in, &n, &l, &r); err != nil {
       return
   }
   px := make([]float64, n)
   py := make([]float64, n)
   cosA := make([]float64, n)
   sinA := make([]float64, n)
   for i := 0; i < n; i++ {
       var x, y, a float64
       fmt.Fscan(in, &x, &y, &a)
       px[i] = x
       py[i] = y
       rad := a * math.Pi / 180.0
       cosA[i] = math.Cos(rad)
       sinA[i] = math.Sin(rad)
   }
   size := 1 << n
   dp := make([]float64, size)
   for i := 0; i < size; i++ {
       dp[i] = l
   }
   for mask := 0; mask < size; mask++ {
       base := dp[mask]
       for i := 0; i < n; i++ {
           if mask>>i&1 == 0 {
               // vector from point to (base,0)
               dx0 := base - px[i]
               dy0 := -py[i]
               // rotate by angle a: (dx0,dy0) * rot
               dirX := dx0*cosA[i] - dy0*sinA[i]
               dirY := dx0*sinA[i] + dy0*cosA[i]
               var np float64
               // if direction goes downward (y < 0), intersects y=0
               if dirY < 0.0 {
                   t := -py[i] / dirY
                   np = px[i] + dirX*t
                   if np > r {
                       np = r
                   }
               } else {
                   np = r
               }
               next := mask | (1 << i)
               if np > dp[next] {
                   dp[next] = np
               }
           }
       }
   }
   res := dp[size-1] - l
   // output fixed with precision 10
   fmt.Printf("%.10f\n", res)
}
