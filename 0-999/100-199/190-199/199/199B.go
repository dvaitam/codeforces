package main

import (
   "fmt"
)

func main() {
   var x1, y1, r1, R1 int64
   var x2, y2, r2, R2 int64
   if _, err := fmt.Scan(&x1, &y1, &r1, &R1); err != nil {
       return
   }
   if _, err := fmt.Scan(&x2, &y2, &r2, &R2); err != nil {
       return
   }
   dx := x1 - x2
   dy := y1 - y2
   d2 := dx*dx + dy*dy
   cnt := 0
   // Outer boundary of ring1
   if (r2 > R1 && d2 < (r2-R1)*(r2-R1)) || (d2 > (R2+R1)*(R2+R1)) {
       cnt++
   }
   // Inner boundary of ring1
   if (r2 > r1 && d2 < (r2-r1)*(r2-r1)) || (d2 > (R2+r1)*(R2+r1)) {
       cnt++
   }
   // Outer boundary of ring2
   if (r1 > R2 && d2 < (r1-R2)*(r1-R2)) || (d2 > (R1+R2)*(R1+R2)) {
       cnt++
   }
   // Inner boundary of ring2
   if (r1 > r2 && d2 < (r1-r2)*(r1-r2)) || (d2 > (R1+r2)*(R1+r2)) {
       cnt++
   }
   fmt.Println(cnt)
}
