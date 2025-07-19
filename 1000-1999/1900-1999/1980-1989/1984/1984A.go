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
       fmt.Fscan(reader, &n)
       arr := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &arr[i])
       }
       if arr[0] == arr[n-1] {
           fmt.Fprintln(writer, "NO")
       } else {
           fmt.Fprintln(writer, "YES")
           // build colors: 'R' at position 1 (0-based index), 'B' elsewhere
           for i := 0; i < n; i++ {
               if i == 1 {
                   writer.WriteByte('R')
               } else {
                   writer.WriteByte('B')
               }
           }
           writer.WriteByte('\n')
       }
   }
}
