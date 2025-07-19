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
       var x, y, n int
       fmt.Fscan(reader, &x, &y, &n)
       b := make([]int, n)
       b[0] = x
       b[n-1] = y
       // reverse slice
       for i := 0; i < n/2; i++ {
           b[i], b[n-1-i] = b[n-1-i], b[i]
       }
       if n > 1 {
           b[1] = b[0] - 1
           for i := 2; i < n-1; i++ {
               b[i] = b[i-1] - i
           }
       }
       // reverse back
       for i := 0; i < n/2; i++ {
           b[i], b[n-1-i] = b[n-1-i], b[i]
       }
       ok := true
       if n > 1 {
           last := b[1] - b[0]
           for i := 2; i < n; i++ {
               diff := b[i] - b[i-1]
               if last <= diff || b[i] == b[i-1] {
                   ok = false
                   break
               }
               last = diff
           }
       }
       if !ok {
           fmt.Fprintln(writer, -1)
       } else {
           for i, v := range b {
               if i > 0 {
                   writer.WriteString(" ")
               }
               fmt.Fprint(writer, v)
           }
           fmt.Fprintln(writer)
       }
   }
}
