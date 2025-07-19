package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

// point represents a 2D point
type point struct { x, y float64 }

// cross returns the cross product of vectors (a->b) x (a->c)
func cross(a, b, c point) float64 {
   return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}

// areaPoly returns the area of a polygon given in order
func areaPoly(p []point) float64 {
   s := 0.0
   n := len(p)
   for i := 1; i < n-1; i++ {
       s += cross(p[0], p[i], p[i+1])
   }
   return math.Abs(s) * 0.5
}

// areaTri returns the area of triangle p1-p2-p3
func areaTri(p1, p2, p3 point) float64 {
   return math.Abs(cross(p1, p2, p3)) * 0.5
}

// convexHull computes the convex hull of pts (removing collinear points) and returns hull in CCW order
func convexHull(pts []point) []point {
   n := len(pts)
   if n <= 1 {
       return append([]point(nil), pts...)
   }
   // sort by y, then x
   sort.Slice(pts, func(i, j int) bool {
       if math.Abs(pts[i].y-pts[j].y) > 1e-9 {
           return pts[i].y < pts[j].y
       }
       return pts[i].x < pts[j].x
   })
   // build lower hull
   lower := make([]point, 0, n)
   for _, p := range pts {
       for len(lower) >= 2 && cross(lower[len(lower)-2], lower[len(lower)-1], p) <= 0 {
           lower = lower[:len(lower)-1]
       }
       lower = append(lower, p)
   }
   // build upper hull
   upper := make([]point, 0, n)
   for i := n - 1; i >= 0; i-- {
       p := pts[i]
       for len(upper) >= 2 && cross(upper[len(upper)-2], upper[len(upper)-1], p) <= 0 {
           upper = upper[:len(upper)-1]
       }
       upper = append(upper, p)
   }
   // concatenate lower and upper, omit last of each (it's duplicate)
   hull := make([]point, 0, len(lower)+len(upper)-2)
   hull = append(hull, lower...)
   if len(upper) > 1 {
       hull = append(hull, upper[1:len(upper)-1]...)
   }
   return hull
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for {
       var n int
       if _, err := fmt.Fscan(reader, &n); err != nil {
           break
       }
       pts := make([]point, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &pts[i].x, &pts[i].y)
       }
       hull := convexHull(pts)
       c := len(hull)
       baseArea := areaPoly(hull)
       // mark hull points
       vis := make([]bool, n)
       for _, hp := range hull {
           for j := 0; j < n; j++ {
               if math.Abs(pts[j].x-hp.x) < 1e-9 && math.Abs(pts[j].y-hp.y) < 1e-9 {
                   vis[j] = true
               }
           }
       }
       ans := -1.0
       // less than 4 hull points: consider interior points
       if c < 4 {
           for i, p := range pts {
               if vis[i] {
                   continue
               }
               // find minimal triangle area to subtract
               minA := math.Inf(1)
               for j := 0; j < c; j++ {
                   k := (j + 1) % c
                   a := areaTri(p, hull[j], hull[k])
                   if a < minA {
                       minA = a
                   }
               }
               if baseArea-minA > ans {
                   ans = baseArea - minA
               }
           }
       } else {
           // choose any 4 hull points
           for i := 0; i < c-3; i++ {
               for j := i + 1; j < c-2; j++ {
                   for k := j + 1; k < c-1; k++ {
                       for l := k + 1; l < c; l++ {
                           quad := []point{hull[i], hull[j], hull[k], hull[l]}
                           a := areaPoly(quad)
                           if a > ans {
                               ans = a
                           }
                       }
                   }
               }
           }
       }
       fmt.Fprintf(writer, "%.9f\n", ans)
   }
}
