package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Point represents a 2D point
type Point struct{ x, y int64 }

// cross returns cross product (b-a)x(c-a)
func cross(a, b, c Point) int64 {
   return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}

// convexHull returns points on convex hull in CCW order, no duplicate endpoints
func convexHull(a []Point) []Point {
   n := len(a)
   if n <= 1 {
       return append([]Point{}, a...)
   }
   sort.Slice(a, func(i, j int) bool {
       if a[i].x != a[j].x {
           return a[i].x < a[j].x
       }
       return a[i].y < a[j].y
   })
   var lo, up []Point
   for _, p := range a {
       for len(lo) >= 2 && cross(lo[len(lo)-2], lo[len(lo)-1], p) <= 0 {
           lo = lo[:len(lo)-1]
       }
       lo = append(lo, p)
   }
   for i := n - 1; i >= 0; i-- {
       p := a[i]
       for len(up) >= 2 && cross(up[len(up)-2], up[len(up)-1], p) <= 0 {
           up = up[:len(up)-1]
       }
       up = append(up, p)
   }
   // remove duplicate endpoints
   lo = lo[:len(lo)-1]
   up = up[:len(up)-1]
   hull := append(lo, up...)
   return hull
}

// pointInConvex checks if q is inside or on border of convex polygon P (CCW)
func pointInConvex(P []Point, q Point) bool {
   n := len(P)
   if n == 0 {
       return false
   }
   // outside range of first wedge
   if cross(P[0], P[1], q) < 0 || cross(P[0], P[n-1], q) > 0 {
       return false
   }
   // binary search sector
   low, high := 1, n-1
   for high-low > 1 {
       mid := (low + high) / 2
       if cross(P[0], P[mid], q) >= 0 {
           low = mid
       } else {
           high = mid
       }
   }
   // check triangle P[0],P[low],P[low+1]
   return cross(P[low], P[low+1], q) >= 0
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var q int
   fmt.Fscan(in, &q)
   var hull []Point
   // read first three insertions
   pts := make([]Point, 0, q)
   for i := 0; i < 3; i++ {
       var t int; var x, y int64
       fmt.Fscan(in, &t, &x, &y)
       pts = append(pts, Point{x, y})
   }
   hull = convexHull(pts)
   // process remaining queries
   for i := 3; i < q; i++ {
       var t int; var x, y int64
       fmt.Fscan(in, &t, &x, &y)
       p := Point{x, y}
       if t == 1 {
           if pointInConvex(hull, p) {
               continue
           }
           // find initial visible edge
           n := len(hull)
           var L, R int
           if cross(hull[0], hull[1], p) < 0 {
               L, R = 0, 1
           } else if cross(hull[0], hull[n-1], p) > 0 {
               L, R = n-1, 0
           } else {
               low, high := 1, n-1
               for high-low > 1 {
                   mid := (low + high) / 2
                   if cross(hull[0], hull[mid], p) > 0 {
                       low = mid
                   } else {
                       high = mid
                   }
               }
               L, R = low, low+1
           }
           // expand visible edges
           // R forward
           for cross(hull[R%len(hull)], hull[(R+1)%len(hull)], p) < 0 {
               R++
           }
           // L backward
           for cross(hull[(L-1+len(hull))%len(hull)], hull[L%len(hull)], p) < 0 {
               L--
           }
           // build new hull starting at p, include tangent vertices at R and L
           newHull := []Point{p}
           n0 := len(hull)
           idx := R % n0
           endIdx := L % n0
           for {
               newHull = append(newHull, hull[idx])
               if idx == endIdx {
                   break
               }
               idx = (idx + 1) % n0
           }
           // rotate so minimal point first
           m := len(newHull)
           minI := 0
           for i := 1; i < m; i++ {
               if newHull[i].x < newHull[minI].x || (newHull[i].x == newHull[minI].x && newHull[i].y < newHull[minI].y) {
                   minI = i
               }
           }
           // rotate
           hull = append(newHull[minI:], newHull[:minI]...)
       } else {
           if pointInConvex(hull, p) {
               fmt.Fprintln(out, "YES")
           } else {
               fmt.Fprintln(out, "NO")
           }
       }
   }
}
