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
       seen := make(map[int]bool)
       dup := false
       for i := 0; i < n; i++ {
           var b int
           fmt.Fscan(reader, &b)
           if seen[b] {
               dup = true
           }
           seen[b] = true
       }
       if dup {
           fmt.Fprintln(writer, "YES")
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}
