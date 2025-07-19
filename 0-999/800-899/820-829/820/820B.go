package main

import (
   "fmt"
   "math"
)

func main() {
   var n, a int64
   if _, err := fmt.Scan(&n, &a); err != nil {
       return
   }
   centerStep := 360.0 / float64(n)
   center := centerStep
   minDiff := math.Abs(center/2.0 - float64(a))
   vC, vS := int64(3), int64(2)
   ver := int64(3)
   center += centerStep
   for center <= 180.0 {
       diff1 := math.Abs(center/2.0 - float64(a))
       if diff1 < minDiff {
           vC = ver + 1
           vS = ver
           minDiff = diff1
       }
       diff2 := math.Abs(180.0 - center/2.0 - float64(a))
       if diff2 < minDiff {
           vC = ver - 1
           vS = ver
           minDiff = diff2
       }
       ver++
       center += centerStep
   }
   fmt.Println(1, vC, vS)
}
