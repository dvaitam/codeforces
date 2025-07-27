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
       var n, k int
       fmt.Fscan(reader, &n, &k)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       // count number of increase events (jumps)
       cntJumps := 0
       for i := 1; i < n; i++ {
           if a[i] > a[i-1] {
               cntJumps++
           }
       }
       // special case: k == 1
       if k == 1 {
           if cntJumps > 0 {
               fmt.Fprintln(writer, -1)
           } else {
               fmt.Fprintln(writer, 1)
           }
           continue
       }
       // each array can have at most k-1 jumps
       capPer := k - 1
       // minimum arrays needed to cover all jumps
       m := (cntJumps + capPer - 1) / capPer
       if m < 1 {
           m = 1
       }
       fmt.Fprintln(writer, m)
   }
}
