package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var N int
   if _, err := fmt.Fscan(reader, &N); err != nil {
       return
   }
   xs := make([]int, N)
   ys := make([]int, N)
   maxX, maxY := 0, 0
   for i := 0; i < N; i++ {
       fmt.Fscan(reader, &xs[i], &ys[i])
       if xs[i] > maxX {
           maxX = xs[i]
       }
       if ys[i] > maxY {
           maxY = ys[i]
       }
   }
   bestArea := 1e30
   // Try A > maxX
   for A := maxX + 1; ; A++ {
       B := 0
       for i := 0; i < N; i++ {
           num := A * ys[i]
           den := A - xs[i]
           b := num / den
           if num%den != 0 {
               b++
           }
           if b > B {
               B = b
           }
       }
       area := float64(A*B) / 2.0
       if area < bestArea {
           bestArea = area
       }
       // If even minimal possible B = maxY gives area >= current best, stop
       if float64(A*maxY)/2.0 > bestArea {
           break
       }
   }
   fmt.Fprintf(writer, "%.10f", bestArea)
}
