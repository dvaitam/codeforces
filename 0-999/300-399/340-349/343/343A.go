package main

import (
   "fmt"
)

func main() {
   var a, b uint64
   if _, err := fmt.Scan(&a, &b); err != nil {
       return
   }
   var cnt uint64
   for a != 1 || b != 1 {
       if a > b {
           // subtract b from a as many times as possible without reaching zero
           k := (a - 1) / b
           cnt += k
           a -= k * b
       } else {
           // subtract a from b similarly
           k := (b - 1) / a
           cnt += k
           b -= k * a
       }
   }
   // include the initial resistor
   cnt++
   fmt.Println(cnt)
}
