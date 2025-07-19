package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
      fmt.Fscan(reader, &a[i])
   }
   minV, maxV := a[0], a[0]
   for _, v := range a {
      if v < minV {
         minV = v
      }
      if v > maxV {
         maxV = v
      }
   }
   var bestT int64
   bestQ := int64(1<<62 - 1)
   for t := minV; t <= maxV; t++ {
      var sum int64
      for _, v := range a {
         diff := v - t
         if diff < 0 {
            diff = -diff
         }
         if diff > 1 {
            sum += diff - 1
         }
      }
      if sum < bestQ {
         bestQ = sum
         bestT = t
      }
   }
   fmt.Println(bestT, bestQ)
}
