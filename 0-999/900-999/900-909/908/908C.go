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
   var r float64
   if _, err := fmt.Fscan(reader, &n, &r); err != nil {
       return
   }
   pushes := make([]float64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &pushes[i])
   }
   ans := make([]float64, n)
   for i := 0; i < n; i++ {
       mx := r
       for j := 0; j < i; j++ {
           dx := math.Abs(pushes[j] - pushes[i])
           if dx <= 2*r {
               t := ans[j] + math.Sqrt(4*r*r - dx*dx)
               if t > mx {
                   mx = t
               }
           }
       }
       ans[i] = mx
   }
   for i := 0; i < n; i++ {
       if i > 0 {
           fmt.Print(" ")
       }
       fmt.Printf("%.10f", ans[i])
   }
   fmt.Println()
}
