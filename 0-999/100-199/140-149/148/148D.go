package main

import (
   "fmt"
   "os"
)

func main() {
   var w, b int
   if _, err := fmt.Fscan(os.Stdin, &w, &b); err != nil {
       return
   }
   var ans, v float64
   cnt := 0
   total := w + b
   for i := 0; cnt < total; i++ {
       if i&1 == 1 {
           cnt++
       } else {
           ans += (1 - v) * float64(w) / float64(total-i)
       }
       cnt++
       v += (1 - v) * float64(w) / float64(total-i)
   }
   fmt.Printf("%.15f", ans)
}
