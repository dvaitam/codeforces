package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   p := make([][2]float64, n)
   for i := 0; i < n; i++ {
       var xi, yi int
       fmt.Fscan(reader, &xi, &yi)
       p[i][0] = float64(xi)
       p[i][1] = float64(yi)
   }
   dpPrev := make([][2]float64, n)
   dpNext := make([][2]float64, n)
   // dpPrev represents dp for length = 1 (initially zero)
   for length := 1; length <= n-1; length++ {
       // reset dpNext
       for i := 0; i < n; i++ {
           dpNext[i][0] = 0
           dpNext[i][1] = 0
       }
       for i := 0; i < n; i++ {
           // state t = 0
           d0 := dpPrev[i][0]
           start := i + 1
           if start >= n {
               start -= n
           }
           end := i - length
           if end < 0 {
               end += n
           }
           distStart := math.Hypot(p[start][0]-p[i][0], p[start][1]-p[i][1])
           if val := d0 + distStart; dpNext[start][0] < val {
               dpNext[start][0] = val
           }
           distEnd := math.Hypot(p[end][0]-p[i][0], p[end][1]-p[i][1])
           if val := d0 + distEnd; dpNext[end][1] < val {
               dpNext[end][1] = val
           }
           // state t = 1
           d1 := dpPrev[i][1]
           end2 := i + length
           if end2 >= n {
               end2 -= n
           }
           start2 := i - 1
           if start2 < 0 {
               start2 += n
           }
           distStart2 := math.Hypot(p[start2][0]-p[i][0], p[start2][1]-p[i][1])
           if val := d1 + distStart2; dpNext[start2][1] < val {
               dpNext[start2][1] = val
           }
           distEnd2 := math.Hypot(p[end2][0]-p[i][0], p[end2][1]-p[i][1])
           if val := d1 + distEnd2; dpNext[end2][0] < val {
               dpNext[end2][0] = val
           }
       }
       // swap dpPrev and dpNext
       dpPrev, dpNext = dpNext, dpPrev
   }
   // find answer in dpPrev (length = n)
   ans := 0.0
   for i := 0; i < n; i++ {
       if dpPrev[i][0] > ans {
           ans = dpPrev[i][0]
       }
       if dpPrev[i][1] > ans {
           ans = dpPrev[i][1]
       }
   }
   fmt.Printf("%.15f", ans)
}
