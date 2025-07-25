package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type point struct{
   x, y int64
}

// cross returns the cross product of OA x OB = (A->O->B)
func cross(a, b, c point) int64 {
   return (b.x - a.x)*(c.y - a.y) - (b.y - a.y)*(c.x - a.x)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   pts := make([]point, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &pts[i].x, &pts[i].y)
   }
   if n == 0 {
       fmt.Println("Ani")
       return
   }
   // sort lex by x, then y
   sort.Slice(pts, func(i, j int) bool {
       if pts[i].x != pts[j].x {
           return pts[i].x < pts[j].x
       }
       return pts[i].y < pts[j].y
   })
   // build convex hull (excluding collinear points)
   var lower []point
   for _, p := range pts {
       for len(lower) >= 2 && cross(lower[len(lower)-2], lower[len(lower)-1], p) <= 0 {
           lower = lower[:len(lower)-1]
       }
       lower = append(lower, p)
   }
   var upper []point
   for i := n-1; i >= 0; i-- {
       p := pts[i]
       for len(upper) >= 2 && cross(upper[len(upper)-2], upper[len(upper)-1], p) <= 0 {
           upper = upper[:len(upper)-1]
       }
       upper = append(upper, p)
   }
   // concatenate lower and upper to get full hull, excluding duplicate endpoints
   hull := lower[:len(lower)-1]
   hull = append(hull, upper[:len(upper)-1]...)
   // check all hull vertices for even parity
   for _, p := range hull {
       if p.x%2 != 0 || p.y%2 != 0 {
           fmt.Println("Ani")
           return
       }
   }
   fmt.Println("Borna")
}
