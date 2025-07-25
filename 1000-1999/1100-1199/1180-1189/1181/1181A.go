package main

import (
   "fmt"
)

func main() {
   var x, y, z int64
   if _, err := fmt.Scan(&x, &y, &z); err != nil {
       return
   }
   a := x / z
   b := y / z
   total := (x + y) / z
   // Determine minimal transfer needed
   var transfer int64
   if total == a+b {
       transfer = 0
   } else {
       r1 := z - x%z
       r2 := z - y%z
       if r1 < r2 {
           transfer = r1
       } else {
           transfer = r2
       }
   }
   fmt.Println(total, transfer)
