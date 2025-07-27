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

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   if n == 1 {
       fmt.Fprintln(writer, a[0])
       return
   }
   if n == 2 {
       fmt.Fprintln(writer, -(a[0] + a[1]))
       return
   }
   var sumAbs int64
   minAbs := a[0]
   if minAbs < 0 {
       minAbs = -minAbs
   }
   for i := 0; i < n; i++ {
       ai := a[i]
       if ai < 0 {
           ai = -ai
       }
       sumAbs += ai
       if ai < minAbs {
           minAbs = ai
       }
   }
   if n%2 == 1 {
       // odd length: one negative inevitable
       fmt.Fprintln(writer, sumAbs - 2*minAbs)
   } else {
       fmt.Fprintln(writer, sumAbs)
   }
}
