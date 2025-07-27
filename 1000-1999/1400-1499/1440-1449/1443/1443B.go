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
   fmt.Fscan(reader, &t)
   for tc := 0; tc < t; tc++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       var s string
       fmt.Fscan(reader, &s)

       cost := 0
       inBlock := false
       seen := false
       zeros := 0
       for i := 0; i < len(s); i++ {
           if s[i] == '1' {
               if !inBlock {
                   if !seen {
                       // first segment: pay activation cost
                       cost += a
                       seen = true
                   } else {
                       // subsequent segment: decide to fill gap or activate separately
                       if zeros*b < a {
                           cost += zeros * b
                       } else {
                           cost += a
                       }
                   }
                   inBlock = true
               }
               // reset zeros counter when in a block
               zeros = 0
           } else {
               if inBlock {
                   // just exited a block, start counting zeros
                   inBlock = false
                   zeros = 1
               } else if seen {
                   // zeros between segments
                   zeros++
               }
           }
       }
       fmt.Fprintln(writer, cost)
   }
}
