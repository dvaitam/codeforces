package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// interval represents a point's coverage [l, r]
type interval struct {
   l, r int64
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   intervals := make([]interval, n)
   for i := 0; i < n; i++ {
       var x, w int64
       fmt.Fscan(in, &x, &w)
       intervals[i].l = x - w
       intervals[i].r = x + w
   }
   sort.Slice(intervals, func(i, j int) bool {
       if intervals[i].r != intervals[j].r {
           return intervals[i].r < intervals[j].r
       }
       return intervals[i].l < intervals[j].l
   })
   count := 0
   var lastR int64 = -1 << 62
   for _, iv := range intervals {
       if iv.l >= lastR {
           count++
           lastR = iv.r
       }
   }
   // output result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, count)
}
