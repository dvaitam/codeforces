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
   for tc := 0; tc < t; tc++ {
       var s string
       fmt.Fscan(reader, &s)
       zeros, ones := 0, 0
       for _, ch := range s {
           if ch == '0' {
               zeros++
           } else {
               ones++
           }
       }
       if zeros == 0 || ones == 0 {
           fmt.Fprint(writer, s)
       } else {
           for zeros > 0 {
               fmt.Fprint(writer, "01")
               zeros--
           }
           for ones > 0 {
               fmt.Fprint(writer, "01")
               ones--
           }
       }
       fmt.Fprint(writer, "\n")
   }
}
