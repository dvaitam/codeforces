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
       xs := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &xs[i])
       }
       diffs := make(map[int]struct{})
       for i := 0; i < n; i++ {
           for j := i + 1; j < n; j++ {
               d := xs[j] - xs[i]
               diffs[d] = struct{}{}
           }
       }
       fmt.Fprintln(writer, len(diffs))
   }
}
