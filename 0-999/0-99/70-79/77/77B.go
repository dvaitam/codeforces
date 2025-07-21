package main

import (
   "bufio"
   "fmt"
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
       var a, b int64
       if _, err := fmt.Fscan(reader, &a, &b); err != nil {
           return
       }
       var prob float64
       if b == 0 {
           // q is always zero, discriminant always >= 0
           prob = 1.0
       } else if a <= 4*b {
           prob = 0.5 + float64(a)/(16.0*float64(b))
       } else {
           prob = 1.0 - float64(b)/float64(a)
       }

       // print with sufficient precision
       fmt.Fprintf(writer, "%.10f\n", prob)
   }
}
