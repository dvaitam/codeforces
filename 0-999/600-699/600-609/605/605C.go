package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

type pkt struct { x, y int64 }

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var p, q int64
   if _, err := fmt.Fscan(in, &n, &p, &q); err != nil {
       return
   }
   a := make([]pkt, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i].x, &a[i].y)
   }
   sort.Slice(a, func(i, j int) bool {
       if a[i].x == a[j].x {
           return a[i].y > a[j].y
       }
       return a[i].x < a[j].x
   })
   // Build convex hull
   s := make([]pkt, 0, n)
   for i := 0; i < n; i++ {
       if i > 0 && a[i].x == a[i-1].x {
           continue
       }
       for len(s) >= 1 && a[i].y >= s[len(s)-1].y {
           s = s[:len(s)-1]
       }
       for len(s) >= 2 {
           sz := len(s)
           x1 := s[sz-1].x - s[sz-2].x
           y1 := s[sz-1].y - s[sz-2].y
           x2 := a[i].x - s[sz-1].x
           y2 := a[i].y - s[sz-1].y
           if x1*y2 - y1*x2 >= 0 {
               s = s[:sz-1]
               continue
           }
           break
       }
       s = append(s, a[i])
   }
   mi := math.Inf(1)
   // Check single points
   for _, pt := range s {
       t := math.Max(float64(p)/float64(pt.x), float64(q)/float64(pt.y))
       if t < mi {
           mi = t
       }
   }
   // Check edges
   for i := 0; i+1 < len(s); i++ {
       a1 := s[i]
       b1 := s[i+1]
       d := float64(a1.x)*float64(b1.y) - float64(a1.y)*float64(b1.x)
       d1 := float64(p)*float64(b1.y) - float64(q)*float64(b1.x)
       d2 := float64(a1.x)*float64(q) - float64(a1.y)*float64(p)
       t := math.Max(0, d1/d) + math.Max(0, d2/d)
       if t < mi {
           mi = t
       }
   }
   out := bufio.NewWriter(os.Stdout)
   fmt.Fprintf(out, "%.10f\n", mi)
   out.Flush()
}
