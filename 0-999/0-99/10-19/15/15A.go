package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, t int
   if _, err := fmt.Fscan(reader, &n, &t); err != nil {
       return
   }
   intervals := make([]struct{ l, r int }, n)
   for i := 0; i < n; i++ {
       var x, a int
       fmt.Fscan(reader, &x, &a)
       intervals[i].l = 2*x - a
       intervals[i].r = 2*x + a
   }
   sort.Slice(intervals, func(i, j int) bool {
       if intervals[i].l != intervals[j].l {
           return intervals[i].l < intervals[j].l
       }
       return intervals[i].r < intervals[j].r
   })
   t2 := 2 * t
   ans := 2
   for i := 0; i+1 < n; i++ {
       gap := intervals[i+1].l - intervals[i].r
       if gap == t2 {
           ans++
       } else if gap > t2 {
           ans += 2
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
}
