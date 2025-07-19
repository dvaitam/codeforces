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
   for ; t > 0; t-- {
       var x int
       fmt.Fscan(reader, &x)
       vc := make([]int, 0, 32)
       for x > 0 {
           if x%2 == 0 {
               vc = append(vc, 0)
               x /= 2
           } else {
               if x%4 == 1 {
                   vc = append(vc, 1)
                   x--
               } else {
                   vc = append(vc, -1)
                   x++
               }
               x /= 2
           }
       }
       fmt.Fprintln(writer, len(vc))
       for _, v := range vc {
           fmt.Fprintf(writer, "%d ", v)
       }
       fmt.Fprintln(writer)
   }
}
