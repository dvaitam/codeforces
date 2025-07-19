package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

var v float64
const eps = 1e-9

func f(c, m, p float64) float64 {
   if math.Abs(p-1) < eps {
       return 1.0
   }
   var e1, e2 float64
   // when c is nearly zero
   if math.Abs(c) < eps {
       mv := math.Min(m, v)
       e2 = f(c, m-mv, p+mv)
       return 1 + m*e2
   }
   // when m is nearly zero
   if math.Abs(m) < eps {
       mv := math.Min(c, v)
       e1 = f(c-mv, m, p+mv)
       return 1 + c*e1
   }
   // general case
   mc := math.Min(c, v)
   mm := math.Min(m, v)
   e1 = f(c-mc, m+mc/2, p+mc/2)
   e2 = f(c+mm/2, m-mm, p+mm/2)
   return 1 + c*e1 + m*e2
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var c, m, p float64
       fmt.Fscan(reader, &c, &m, &p, &v)
       res := f(c, m, p)
       fmt.Fprintf(writer, "%.12f\n", res)
   }
}
