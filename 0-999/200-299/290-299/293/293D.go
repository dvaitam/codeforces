package main

import (
   "bufio"
   "fmt"
   "os"
)

const ZERO = 1100000

var fromArr, toArr []int
var x, y []int

func det(a, b, c, d int64) int64 {
   return a*d - b*c
}

func calc(n int) float64 {
   var area int64
   // close polygon
   x[n] = x[0]
   y[n] = y[0]
   for i := 0; i < n; i++ {
       area += det(int64(x[i]), int64(y[i]), int64(x[i+1]), int64(y[i+1]))
   }
   if area > 0 {
       // reverse vertices
       for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
           x[i], x[j] = x[j], x[i]
           y[i], y[j] = y[j], y[i]
       }
       x[n] = x[0]
       y[n] = y[0]
   }
   // determine scan range
   minX, maxX := 2*ZERO, 0
   for i := 0; i < n; i++ {
       if x[i] < minX {
           minX = x[i]
       }
       if x[i] > maxX {
           maxX = x[i]
       }
   }
   // clear arrays in range
   for tx := minX; tx <= maxX; tx++ {
       fromArr[tx] = 0
       toArr[tx] = 0
   }
   // fill boundaries
   for i := 0; i < n; i++ {
       xi, yi := x[i], y[i]
       xj, yj := x[i+1], y[i+1]
       if xi < xj {
           dx := xj - xi
           for tx := xi; tx <= xj; tx++ {
               ty := (yj*(tx-xi) + yi*(xj-tx)) / dx
               toArr[tx] = ty
           }
       } else if xi > xj {
           dx := xi - xj
           for tx := xj; tx <= xi; tx++ {
               ty := (yi*(tx-xj) + yj*(xi-tx) + dx - 1) / dx
               fromArr[tx] = ty
           }
       } else {
           // vertical segment
           y0, y1 := yi, yj
           if y0 > y1 {
               y0, y1 = y1, y0
           }
           fromArr[xi] = y0
           toArr[xi] = y1
       }
   }
   // accumulate answer
   var ans, sum, sumSq, cnt float64
   for tx := minX; tx <= maxX; tx++ {
       cur := toArr[tx] - fromArr[tx] + 1
       if cur <= 0 {
           continue
       }
       fcur := float64(cur)
       ftx := float64(tx)
       ans += fcur*cnt*ftx*ftx + fcur*sumSq - 2*fcur*ftx*sum
       sum += fcur * ftx
       sumSq += fcur * ftx * ftx
       cnt += fcur
   }
   return ans * 2 / (cnt * (cnt - 1))
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   x = make([]int, n+1)
   y = make([]int, n+1)
   fromArr = make([]int, 2*ZERO+1)
   toArr = make([]int, 2*ZERO+1)
   for i := 0; i < n; i++ {
       var xi, yi int
       fmt.Fscan(reader, &xi, &yi)
       x[i] = xi + ZERO
       y[i] = yi + ZERO
   }
   ans := calc(n)
   // swap axes and recalc
   for i := 0; i < n; i++ {
       x[i], y[i] = y[i], x[i]
   }
   ans += calc(n)
   fmt.Printf("%.6f\n", ans/2)
}
