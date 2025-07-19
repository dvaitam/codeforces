package main

import (
   "fmt"
)

func main() {
   var s1, s2 string
   if _, err := fmt.Scan(&s1, &s2); err != nil {
       return
   }
   var h1, m1, h2, m2 int
   // parse times in format hh:mm
   fmt.Sscanf(s1, "%d:%d", &h1, &m1)
   fmt.Sscanf(s2, "%d:%d", &h2, &m2)
   t1 := h1*60 + m1
   t2 := h2*60 + m2
   mid := (t1 + t2) / 2
   h := mid / 60
   m := mid % 60
   fmt.Printf("%02d:%02d", h, m)
}
