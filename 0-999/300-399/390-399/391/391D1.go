package main

import (
   "bufio"
   "fmt"
   "os"
)

func min4(a, b, c, d int64) int64 {
   m := a
   if b < m {
       m = b
   }
   if c < m {
       m = c
   }
   if d < m {
       m = d
   }
   return m
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   type SegV struct{ x, y1, y2 int64 }
   type SegH struct{ y, x1, x2 int64 }
   vs := make([]SegV, n)
   for i := 0; i < n; i++ {
       var xi, yi, li int64
       fmt.Fscan(reader, &xi, &yi, &li)
       vs[i] = SegV{xi, yi, yi + li}
   }
   hs := make([]SegH, m)
   for i := 0; i < m; i++ {
       var xi, yi, li int64
       fmt.Fscan(reader, &xi, &yi, &li)
       hs[i] = SegH{yi, xi, xi + li}
   }
   var best int64 = 0
   for _, v := range vs {
       for _, h := range hs {
           // check intersection
           if h.x1 <= v.x && v.x <= h.x2 && v.y1 <= h.y && h.y <= v.y2 {
               d1 := h.y - v.y1
               d2 := v.y2 - h.y
               d3 := v.x - h.x1
               d4 := h.x2 - v.x
               size := min4(d1, d2, d3, d4)
               if size > best {
                   best = size
               }
           }
       }
   }
   fmt.Println(best)
}
