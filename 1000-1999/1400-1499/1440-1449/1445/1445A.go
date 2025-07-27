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
       var n, x int
       fmt.Fscan(reader, &n, &x)
       a := make([]int, n)
       b := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &b[i])
       }
       ok := true
       // a and b are given sorted ascending
       for i := 0; i < n; i++ {
           if a[i]+b[n-1-i] > x {
               ok = false
               break
           }
       }
       if ok {
           fmt.Fprintln(writer, "Yes")
       } else {
           fmt.Fprintln(writer, "No")
       }
   }
}
