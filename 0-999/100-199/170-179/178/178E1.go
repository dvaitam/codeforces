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
   fmt.Fscan(reader, &n)
   grid := make([]byte, n*n)
   for i := 0; i < n*n; i++ {
       var v int
       fmt.Fscan(reader, &v)
       if v != 0 {
           grid[i] = 1
       }
   }
   visited := make([]bool, n*n)
   circles := 0
   squares := 0
   // 4-connectivity directions: up, down, left, right
   dirs := []int{-n, n, -1, 1}
   for idx := 0; idx < n*n; idx++ {
       if grid[idx] == 1 && !visited[idx] {
           // BFS for component
           q := []int{idx}
           visited[idx] = true
           sumX, sumY, sumXXYY := 0.0, 0.0, 0.0
           area := 0
           for qi := 0; qi < len(q); qi++ {
               u := q[qi]
               x := u / n
               y := u % n
               area++
               fx := float64(x)
               fy := float64(y)
               sumX += fx
               sumY += fy
               sumXXYY += fx*fx + fy*fy
               for _, d := range dirs {
                   v := u + d
                   // check horizontal wrap
                   if d == -1 && y == 0 {
                       continue
                   }
                   if d == 1 && y == n-1 {
                       continue
                   }
                   if v < 0 || v >= len(grid) {
                       continue
                   }
                   if grid[v] == 1 && !visited[v] {
                       visited[v] = true
                       q = append(q, v)
                   }
               }
           }
           // ignore small noise components
           if area < 100 {
               continue
           }
           // compute centroid and second moment
           areaF := float64(area)
           xc := sumX / areaF
           yc := sumY / areaF
           E := sumXXYY/areaF - (xc*xc + yc*yc)
           // normalized moment for circle vs square
           Nc := 2 * math.Pi * E / areaF
           // ideal Nc for square = π/3 ≈1.047
           if math.Abs(Nc-1) < math.Abs(Nc-math.Pi/3) {
               circles++
           } else {
               squares++
           }
       }
   }
   fmt.Fprintln(writer, circles, squares)
}
