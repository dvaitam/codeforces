package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   rd := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(rd, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var d float64
       fmt.Fscan(rd, &d)
       D := d*d - 4*d
       if D < 0 {
           fmt.Println("N")
       } else {
           sqrtD := math.Sqrt(D)
           a := (d + sqrtD) / 2.0
           b := (d - sqrtD) / 2.0
           fmt.Printf("Y %.15f %.15f\n", a, b)
       }
   }
}
