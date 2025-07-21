package main

import (
   "fmt"
   "math"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   // Read brick projections
   x1 := make([]float64, n+1)
   y1 := make([]float64, n+1)
   x2 := make([]float64, n+1)
   y2 := make([]float64, n+1)
   minX := make([]float64, n+1)
   maxX := make([]float64, n+1)
   minY := make([]float64, n+1)
   maxY := make([]float64, n+1)
   cx := make([]float64, n+1)
   cy := make([]float64, n+1)
   w := make([]float64, n+1)
   for i := 1; i <= n; i++ {
       var xi1, yi1, xi2, yi2 int
       fmt.Scan(&xi1, &yi1, &xi2, &yi2)
       x1[i], y1[i] = float64(xi1), float64(yi1)
       x2[i], y2[i] = float64(xi2), float64(yi2)
       // compute bounds
       minX[i] = math.Min(x1[i], x2[i])
       maxX[i] = math.Max(x1[i], x2[i])
       minY[i] = math.Min(y1[i], y2[i])
       maxY[i] = math.Max(y1[i], y2[i])
       // center and weight
       cx[i] = (x1[i] + x2[i]) / 2.0
       cy[i] = (y1[i] + y2[i]) / 2.0
       a := maxX[i] - minX[i]
       w[i] = a * a * a
   }
   const eps = 1e-9
   // try building tower
   for k := 2; k <= n; k++ {
       // check stability after adding brick k
       for i := k; i >= 2; i-- {
           // compute center of mass of bricks i..k
           var sumW, sumX, sumY float64
           for j := i; j <= k; j++ {
               sumW += w[j]
               sumX += w[j] * cx[j]
               sumY += w[j] * cy[j]
           }
           cmX := sumX / sumW
           cmY := sumY / sumW
           // support area between brick i-1 and i
           sx := math.Max(minX[i-1], minX[i])
           tx := math.Min(maxX[i-1], maxX[i])
           sy := math.Max(minY[i-1], minY[i])
           ty := math.Min(maxY[i-1], maxY[i])
           if cmX < sx-eps || cmX > tx+eps || cmY < sy-eps || cmY > ty+eps {
               // unstable when adding brick k, so maximal stable tower has k-1 bricks
               fmt.Println(k - 1)
               return
           }
       }
   }
   // all stable
   fmt.Println(n)
}
