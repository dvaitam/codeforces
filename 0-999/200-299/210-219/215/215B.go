package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   // read n and x values to choose maximum r1
   fmt.Fscan(reader, &n)
   r1 := 0
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if x > r1 {
           r1 = x
       }
   }
   // read m and y values to choose maximum p1
   fmt.Fscan(reader, &n)
   p1 := 0
   for i := 0; i < n; i++ {
       var y int
       fmt.Fscan(reader, &y)
       if y > p1 {
           p1 = y
       }
   }
   // read k and z values to choose minimum p2
   fmt.Fscan(reader, &n)
   p2 := int(1<<31 - 1)
   for i := 0; i < n; i++ {
       var z int
       fmt.Fscan(reader, &z)
       if z < p2 {
           p2 = z
       }
   }
   // read constants A and B
   var A, B int
   fmt.Fscan(reader, &A, &B)

   // compute r2 = sqrt( B*p1*r1^2 / (A*p2 + B*p1) )
   num := float64(B * p1 * r1 * r1)
   den := float64(A*p2 + B*p1)
   r2 := math.Sqrt(num / den)
   fmt.Printf("%.10f\n", r2)
}
