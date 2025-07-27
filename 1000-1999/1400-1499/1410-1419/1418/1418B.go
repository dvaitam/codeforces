package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       locked := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &locked[i])
       }
       // collect unlocked values
       var free []int
       for i := 0; i < n; i++ {
           if locked[i] == 0 {
               free = append(free, a[i])
           }
       }
       // sort descending to maximize prefix sums
       sort.Sort(sort.Reverse(sort.IntSlice(free)))
       // assign back
       idx := 0
       for i := 0; i < n; i++ {
           if locked[i] == 0 {
               a[i] = free[idx]
               idx++
           }
       }
       // output
       for i, v := range a {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprintf(writer, "%d", v)
       }
       writer.WriteByte('\n')
   }
}
