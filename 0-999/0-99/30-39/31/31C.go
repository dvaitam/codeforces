package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   type interval struct { l, r, idx int }
   intervals := make([]interval, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &intervals[i].l, &intervals[i].r)
       intervals[i].idx = i + 1
   }
   sort.Slice(intervals, func(i, j int) bool {
       if intervals[i].l == intervals[j].l {
           return intervals[i].r < intervals[j].r
       }
       return intervals[i].l < intervals[j].l
   })
   // conflict count per original index
   conflict := make([]int, n+1)
   total := 0
   for i := 0; i < n; i++ {
       li := intervals[i].l
       ri := intervals[i].r
       for j := i + 1; j < n; j++ {
           if intervals[j].l >= ri {
               break
           }
           // intervals[j].l < ri => overlap
           conflict[intervals[i].idx]++
           conflict[intervals[j].idx]++
           total++
       }
   }
   // find valid removals
   res := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if conflict[i] == total {
           res = append(res, i)
       }
   }
   // output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, len(res))
   if len(res) > 0 {
       for i, v := range res {
           if i > 0 {
               w.WriteString(" ")
           }
           fmt.Fprint(w, v)
       }
       w.WriteByte('\n')
   }
}
