package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   // Read 4 segments
   type seg struct{ x1, y1, x2, y2 int64 }
   segs := make([]seg, 4)
   for i := 0; i < 4; i++ {
       if _, err := fmt.Fscan(in, &segs[i].x1, &segs[i].y1, &segs[i].x2, &segs[i].y2); err != nil {
           fmt.Println("NO")
           return
       }
   }
   var hor, ver []seg
   for _, s := range segs {
       if s.y1 == s.y2 {
           hor = append(hor, s)
       } else if s.x1 == s.x2 {
           ver = append(ver, s)
       } else {
           fmt.Println("NO")
           return
       }
   }
   if len(hor) != 2 || len(ver) != 2 {
       fmt.Println("NO")
       return
   }
   // Process horizontals
   // get y coordinates
   y0 := hor[0].y1
   y1 := hor[1].y1
   if y0 == y1 {
       fmt.Println("NO")
       return
   }
   // determine y_min, y_max
   var yMin, yMax int64
   if y0 < y1 {
       yMin, yMax = y0, y1
   } else {
       yMin, yMax = y1, y0
   }
   // get x interval for horizontals
   x0l, x0h := min(hor[0].x1, hor[0].x2), max(hor[0].x1, hor[0].x2)
   x1l, x1h := min(hor[1].x1, hor[1].x2), max(hor[1].x1, hor[1].x2)
   if x0l != x1l || x0h != x1h {
       fmt.Println("NO")
       return
   }
   if x0l >= x0h {
       fmt.Println("NO")
       return
   }
   xMin, xMax := x0l, x0h

   // Process verticals
   xv0 := ver[0].x1
   xv1 := ver[1].x1
   if xv0 == xv1 {
       fmt.Println("NO")
       return
   }
   var xvMin, xvMax int64
   if xv0 < xv1 {
       xvMin, xvMax = xv0, xv1
   } else {
       xvMin, xvMax = xv1, xv0
   }
   if xvMin != xMin || xvMax != xMax {
       fmt.Println("NO")
       return
   }
   y0l, y0h := min(ver[0].y1, ver[0].y2), max(ver[0].y1, ver[0].y2)
   y1l, y1h := min(ver[1].y1, ver[1].y2), max(ver[1].y1, ver[1].y2)
   if y0l != y1l || y0h != y1h {
       fmt.Println("NO")
       return
   }
   if y0l != yMin || y0h != yMax {
       fmt.Println("NO")
       return
   }
   // All checks passed
   fmt.Println("YES")
}

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}
