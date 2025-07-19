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

   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for tc := 0; tc < T; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       var lastNeg int
       hasNeg := false
       var maxVal int
       for i := 0; i < n; i++ {
           var x int
           fmt.Fscan(reader, &x)
           if x < 0 {
               hasNeg = true
               lastNeg = x
           }
           if i == 0 || x > maxVal {
               maxVal = x
           }
       }
       if hasNeg {
           fmt.Fprintln(writer, lastNeg)
       } else {
           fmt.Fprintln(writer, maxVal)
       }
   }
}
