package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   var s, t string
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   if _, err := fmt.Scan(&t); err != nil {
       return
   }
   total := 0
   for i := 0; i < n; i++ {
       a := int(s[i] - '0')
       b := int(t[i] - '0')
       diff := a - b
       if diff < 0 {
           diff = -diff
       }
       if diff > 5 {
           diff = 10 - diff
       }
       total += diff
   }
   fmt.Println(total)
}
