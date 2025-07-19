package main

import (
   "fmt"
)

func main() {
   var n int
   var t int64
   if _, err := fmt.Scan(&n, &t); err != nil {
       return
   }
   bestTime := int64(1 << 62)
   bestIndex := 1
   for i := 1; i <= n; i++ {
       var d, s int64
       fmt.Scan(&d, &s)
       if d < t {
           delta := t - d
           k := (delta + s - 1) / s
           d += k * s
       }
       if d <= bestTime {
           bestTime = d
           bestIndex = i
       }
   }
   fmt.Println(bestIndex)
}
