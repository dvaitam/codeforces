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
       var n int
       var k int64
       fmt.Fscan(reader, &n, &k)
       if n == 1 {
           fmt.Fprintln(writer, k)
           continue
       }
       var x int64 = 0
       for i := 0; ; i++ {
           bit := int64(1) << i
           if x|bit > k {
               break
           }
           x |= bit
       }
       // Prepare answer: first x, then k-x, then zeros
       // Print all n numbers
       // a1 = x, a2 = k-x
       fmt.Fprint(writer, x)
       fmt.Fprint(writer, " ")
       fmt.Fprint(writer, k-x)
       for i := 2; i < n; i++ {
           fmt.Fprint(writer, " 0")
       }
       fmt.Fprint(writer, '\n')
   }
}
