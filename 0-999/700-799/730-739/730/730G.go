package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   type interval struct{ l, r int64 }
   intervals := make([]interval, 0, n)
   for i := 0; i < n; i++ {
       var s, d int64
       fmt.Fscan(in, &s, &d)
       // try preferred start
       var x int64
       sEnd := s + d - 1
       ok := true
       for _, inv := range intervals {
           if !(inv.r < s || inv.l > sEnd) {
               ok = false
               break
           }
       }
       if ok {
           x = s
       } else {
           // find earliest gap of length d
           var prev int64 = 0
           for _, inv := range intervals {
               gapL := prev + 1
               gapR := inv.l - 1
               if gapR >= gapL && gapR-gapL+1 >= d {
                   x = gapL
                   break
               }
               if inv.r > prev {
                   prev = inv.r
               }
           }
           if x == 0 {
               x = prev + 1
           }
       }
       e := x + d - 1
       // output
       fmt.Fprintf(out, "%d %d\n", x, e)
       // insert interval keeping sort by l
       pos := len(intervals)
       for j, inv := range intervals {
           if x < inv.l {
               pos = j
               break
           }
       }
       // extend slice
       intervals = append(intervals, interval{})
       copy(intervals[pos+1:], intervals[pos:])
       intervals[pos] = interval{l: x, r: e}
   }
}
