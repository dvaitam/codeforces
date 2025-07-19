package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var tests int
   if _, err := fmt.Fscan(reader, &tests); err != nil {
       return
   }
   for t := 0; t < tests; t++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       var sb strings.Builder
       // build alternating bits while both a and b are positive
       for a > 0 && b > 0 {
           if a > b {
               sb.WriteByte('0')
               sb.WriteByte('1')
           } else {
               sb.WriteByte('1')
               sb.WriteByte('0')
           }
           a--
           b--
       }
       // append remaining zeros or ones
       for a > 0 {
           sb.WriteByte('0')
           a--
       }
       for b > 0 {
           sb.WriteByte('1')
           b--
       }
       fmt.Fprintln(writer, sb.String())
   }
}
