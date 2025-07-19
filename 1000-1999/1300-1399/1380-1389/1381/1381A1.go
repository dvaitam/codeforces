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

   var tc int
   fmt.Fscan(reader, &tc)
   for tc > 0 {
       tc--
       var n int
       fmt.Fscan(reader, &n)
       var s, tStr string
       fmt.Fscan(reader, &s, &tStr)

       ops := make([]int, 0, 3*n)
       // If first characters differ, flip position 1
       if s[0] != tStr[0] {
           ops = append(ops, 1)
       }
       // For each position i (1-based i+1), if s[i]!=t[i], apply three ops
       for i := 1; i < n; i++ {
           if s[i] == tStr[i] {
               continue
           }
           ops = append(ops, i+1)
           ops = append(ops, 1)
           ops = append(ops, i+1)
       }
       k := len(ops)
       if k == 0 {
           fmt.Fprintln(writer, 0)
       } else {
           // print count and operations
           fmt.Fprint(writer, k)
           fmt.Fprint(writer, " ")
           for i, v := range ops {
               if i > 0 {
                   fmt.Fprint(writer, " ")
               }
               fmt.Fprint(writer, v)
           }
           fmt.Fprintln(writer)
       }
   }
}
