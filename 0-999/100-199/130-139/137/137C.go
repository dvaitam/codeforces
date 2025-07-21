package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   intervals := make([][2]int, n)
   for i := 0; i < n; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       intervals[i][0] = a
       intervals[i][1] = b
   }
   // Sort by start time ascending
   sort.Slice(intervals, func(i, j int) bool {
       return intervals[i][0] < intervals[j][0]
   })
   count := 0
   maxEnd := -1 << 60
   for _, iv := range intervals {
       a, b := iv[0], iv[1]
       if b < maxEnd {
           count++
       } else if b > maxEnd {
           maxEnd = b
       }
   }
   fmt.Fprintln(out, count)
}
