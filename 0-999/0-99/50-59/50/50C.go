package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Point represents a coordinate
type Point struct {
   x, y int64
}

// cross returns cross product of AB x AC
func cross(a, b, c Point) int64 {
   return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}

// abs64 returns absolute value of v
func abs64(v int64) int64 {
   if v < 0 {
       return -v
   }
   return v
}

// convexHull computes the convex hull of pts (must be sorted by x, then y)
// returns points in CCW order, without repeating the first point
func convexHull(pts []Point) []Point {
   n := len(pts)
   if n <= 1 {
       return pts
   }
   var lower []Point
   for _, p := range pts {
       for len(lower) >= 2 && cross(lower[len(lower)-2], lower[len(lower)-1], p) <= 0 {
           lower = lower[:len(lower)-1]
       }
       lower = append(lower, p)
   }
   var upper []Point
   for i := n - 1; i >= 0; i-- {
       p := pts[i]
       for len(upper) >= 2 && cross(upper[len(upper)-2], upper[len(upper)-1], p) <= 0 {
           upper = upper[:len(upper)-1]
       }
       upper = append(upper, p)
   }
   // Concatenate lower and upper, excluding duplicate endpoints
   hull := lower[:len(lower)-1]
   hull = append(hull, upper[:len(upper)-1]...)
   return hull
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   pts := make([]Point, 0, n)
   for i := 0; i < n; i++ {
       var xi, yi int64
       fmt.Fscan(reader, &xi, &yi)
       pts = append(pts, Point{xi, yi})
   }
   // remove duplicates
   sort.Slice(pts, func(i, j int) bool {
       if pts[i].x != pts[j].x {
           return pts[i].x < pts[j].x
       }
       return pts[i].y < pts[j].y
   })
   uniq := pts[:0]
   for i, p := range pts {
       if i == 0 || p.x != pts[i-1].x || p.y != pts[i-1].y {
           uniq = append(uniq, p)
       }
   }
   pts = uniq
   m := len(pts)
   if m == 0 {
       fmt.Fprintln(writer, 0)
       return
   }
   if m == 1 {
       fmt.Fprintln(writer, 4)
       return
   }
   // compute hull
   hull := convexHull(pts)
   hsz := len(hull)
   if hsz == 1 {
       fmt.Fprintln(writer, 4)
       return
   }
   if hsz == 2 {
       dx := abs64(hull[1].x - hull[0].x)
       dy := abs64(hull[1].y - hull[0].y)
       d := dx
       if dy > d {
           d = dy
       }
       // loop around segment
       fmt.Fprintln(writer, 2*d+4)
       return
   }
   // perimeter in Chebyshev (king moves)
   var perim int64
   for i := 0; i < hsz; i++ {
       j := (i + 1) % hsz
       dx := abs64(hull[j].x - hull[i].x)
       dy := abs64(hull[j].y - hull[i].y)
       if dx > dy {
           perim += dx
       } else {
           perim += dy
       }
   }
   fmt.Fprintln(writer, perim)
}
