package main

import (
   "bufio"
   "fmt"
   "os"
)

// Point represents a 2D point with integer coordinates.
type Point struct {
   x, y int64
}

// cross returns the cross product of vectors (b - a) and (c - a).
func cross(a, b, c Point) int64 {
   return (b.x - a.x)*(c.y - a.y) - (b.y - a.y)*(c.x - a.x)
}

// pointInConvex checks if point q is strictly inside the convex polygon poly.
// poly must be in counter-clockwise order and have no three collinear points.
func pointInConvex(poly []Point, q Point) bool {
   n := len(poly)
   // q must be to the left of edge [poly[0]->poly[1]] and to the right of [poly[0]->poly[n-1]]
   if cross(poly[0], poly[1], q) <= 0 {
       return false
   }
   if cross(poly[0], poly[n-1], q) >= 0 {
       return false
   }
   // binary search to find i such that q is in triangle (poly[0], poly[i], poly[i+1])
   l, r := 1, n-1
   for l+1 < r {
       m := (l + r) >> 1
       if cross(poly[0], poly[m], q) > 0 {
           l = m
       } else {
           r = m
       }
   }
   // check q is to the left of edge [poly[l]->poly[l+1]]
   if cross(poly[l], poly[l+1], q) <= 0 {
       return false
   }
   return true
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   poly := make([]Point, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &poly[i].x, &poly[i].y)
   }
   // reverse to convert from clockwise to counter-clockwise
   for i := 0; i < n/2; i++ {
       poly[i], poly[n-1-i] = poly[n-1-i], poly[i]
   }

   var m int
   fmt.Fscan(reader, &m)
   // check each vertex of B is strictly inside A
   for i := 0; i < m; i++ {
       var q Point
       fmt.Fscan(reader, &q.x, &q.y)
       if !pointInConvex(poly, q) {
           fmt.Fprintln(writer, "NO")
           return
       }
   }
   fmt.Fprintln(writer, "YES")
}
