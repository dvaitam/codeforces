package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var r int64
   if _, err := fmt.Fscan(in, &r); err != nil {
       return
   }
   // H(x,y) = x^2 + 2xy + x + 1 = r
   // Solve for y: y = (r - (x^2 + x + 1)) / (2x)
   found := false
   var ansX, ansY int64
   for x := int64(1); ; x++ {
       // minimal value t = x^2 + x + 1; numerator N = r - t
       t := x*x + x + 1
       if t >= r {
           break
       }
       N := r - t
       denom := 2 * x
       if N%denom != 0 {
           continue
       }
       y := N / denom
       if y >= 1 {
           ansX, ansY = x, y
           found = true
           break
       }
   }
   if found {
       fmt.Printf("%d %d", ansX, ansY)
   } else {
       fmt.Println("NO")
   }
}
