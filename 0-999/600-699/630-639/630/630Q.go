package main

import (
   "fmt"
   "math"
)

func main() {
   var ret float64
   for i := 3; i <= 5; i++ {
       var length float64
       if _, err := fmt.Scan(&length); err != nil {
           return
       }
       // Circumradius of regular i-gon with side length
       r := length / 2 / math.Sin(math.Pi/float64(i))
       // Height of pyramid side
       h := math.Sqrt(length*length - r*r)
       // Area of regular i-gon
       area := r*r*math.Sin(2*math.Pi/float64(i))/2 * float64(i)
       // Volume contribution
       ret += area * h / 3
   }
   fmt.Printf("%.10f\n", ret)
}
