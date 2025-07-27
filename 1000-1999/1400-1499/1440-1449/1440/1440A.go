package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for tc := 0; tc < t; tc++ {
       var n, c0, c1, h int
       fmt.Fscan(reader, &n, &c0, &c1, &h)
       var s string
       fmt.Fscan(reader, &s)
       count0, count1 := 0, 0
       for i := range s {
           if s[i] == '0' {
               count0++
           } else {
               count1++
           }
       }
       // Determine optimal conversion costs
       if c0 > c1+h {
           c0 = c1 + h
       }
       if c1 > c0+h {
           c1 = c0 + h
       }
       result := count0*c0 + count1*c1
       fmt.Fprintln(writer, result)
   }
}
