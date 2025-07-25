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

   var q int
   if _, err := fmt.Fscan(reader, &q); err != nil {
       return
   }
   for qi := 0; qi < q; qi++ {
       var n int
       fmt.Fscan(reader, &n)
       xMin, xMax := -100000, 100000
       yMin, yMax := -100000, 100000
       for i := 0; i < n; i++ {
           var x, y int
           var f1, f2, f3, f4 int
           fmt.Fscan(reader, &x, &y, &f1, &f2, &f3, &f4)
           if f1 == 0 && x > xMin {
               xMin = x
           }
           if f3 == 0 && x < xMax {
               xMax = x
           }
           if f4 == 0 && y > yMin {
               yMin = y
           }
           if f2 == 0 && y < yMax {
               yMax = y
           }
       }
       if xMin <= xMax && yMin <= yMax {
           fmt.Fprintln(writer, 1, xMin, yMin)
       } else {
           fmt.Fprintln(writer, 0)
       }
   }
}
