package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type point struct {
   x, y int64
}

func quadrant(p point) int {
   if p.x > 0 && p.y >= 0 {
       return 1
   } else if p.x <= 0 && p.y > 0 {
       return 2
   } else if p.x < 0 && p.y <= 0 {
       return 3
   }
   return 4
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var v []point
   var sx, sy int64
   for i := 0; i < n; i++ {
       var a, b, c, d int64
       fmt.Fscan(reader, &a, &b, &c, &d)
       dx := a - b
       dy := c - d
       if dx == 0 && dy == 0 {
           continue
       }
       v = append(v, point{dx, dy})
       v = append(v, point{-dx, -dy})
       if c < d || (c == d && a < b) {
           sx += dx
           sy += dy
       }
   }
   sort.Slice(v, func(i, j int) bool {
       pi, pj := v[i], v[j]
       qi, qj := quadrant(pi), quadrant(pj)
       if qi != qj {
           return qi < qj
       }
       // sort by cross product descending: pi.x* pj.y > pi.y*pj.x
       return pi.x*pj.y > pi.y*pj.x
   })
   var ans float64
   // initial squared magnitude
   ans = float64(sx*sx + sy*sy)
   for _, p := range v {
       sx += p.x
       sy += p.y
       val := float64(sx*sx + sy*sy)
       if val > ans {
           ans = val
       }
   }
   fmt.Fprintf(writer, "%.10f\n", ans)
}
