package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var d, k int64
       fmt.Fscan(reader, &d, &k)
       // D2 = d^2, k2 = k^2
       D2 := d * d
       k2 := k * k
       // compute x = max x >=0: 2*(x*k)^2 <= d^2
       // i.e., 2*x^2*k^2 <= D2
       // initial estimate
       x := int64(math.Sqrt(float64(D2) / float64(2*k2)))
       // adjust x
       for (x+1)*(x+1)*2*k2 <= D2 {
           x++
       }
       for x >= 0 && x*x*2*k2 > D2 {
           x--
       }
       // remainder after x moves on both axes
       rem := D2 - x*x*k2
       // compute y = floor(sqrt(rem) / k)
       y := int64(math.Sqrt(float64(rem) / float64(k2)))
       // adjust y
       for (y+1)*(y+1)*k2 <= rem {
           y++
       }
       for y >= 0 && y*y*k2 > rem {
           y--
       }
       // if y > x, odd moves available -> first wins
       if y > x {
           fmt.Fprintln(writer, "Ashish")
       } else {
           fmt.Fprintln(writer, "Utkarsh")
       }
   }
}
