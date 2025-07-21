package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

type Point struct{ x, y int64 }

// pointOnSegment checks if p is on segment a-b
func pointOnSegment(a, b, p Point) bool {
   // cross product and bounding box
   dx1 := b.x - a.x
   dy1 := b.y - a.y
   dx2 := p.x - a.x
   dy2 := p.y - a.y
   if dx1*dy2 != dy1*dx2 {
       return false
   }
   if dx2 < 0 && dx1 > 0 || dx2 > 0 && dx1 < 0 {
       return false
   }
   if dy2 < 0 && dy1 > 0 || dy2 > 0 && dy1 < 0 {
       return false
   }
   if abs64(dx2) > abs64(dx1) || abs64(dy2) > abs64(dy1) {
       return false
   }
   return true
}

func abs64(a int64) int64 {
   if a < 0 {
       return -a
   }
   return a
}

// pointInPoly returns true if p is inside or on boundary of poly
func pointInPoly(poly []Point, p Point) bool {
   inside := false
   n := len(poly)
   for i := 0; i < n; i++ {
       a := poly[i]
       b := poly[(i+1)%n]
       if pointOnSegment(a, b, p) {
           return true
       }
       // ray casting to the right
       yi, yj := a.y, b.y
       xi, xj := a.x, b.x
       if (yi > p.y) != (yj > p.y) {
           // compute x coordinate of intersection
           x := xi + (p.y-yi)*(xj-xi)/(yj-yi)
           if x >= p.x {
               inside = !inside
           }
       }
   }
   return inside
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, _ := reader.ReadString('\n')
   parts := strings.Fields(line)
   n, _ := strconv.Atoi(parts[0])
   m, _ := strconv.Atoi(parts[1])
   // skip edges
   for i := 0; i < m; i++ {
       reader.ReadString('\n')
   }
   pts := make([]Point, n)
   for i := 0; i < n; i++ {
       line, _ = reader.ReadString('\n')
       parts = strings.Fields(line)
       x, _ := strconv.ParseInt(parts[0], 10, 64)
       y, _ := strconv.ParseInt(parts[1], 10, 64)
       pts[i] = Point{x, y}
   }
   line, _ = reader.ReadString('\n')
   q, _ := strconv.Atoi(strings.TrimSpace(line))
   out := make([]int, q)
   for qi := 0; qi < q; qi++ {
       line, _ = reader.ReadString('\n')
       parts = strings.Fields(line)
       k, _ := strconv.Atoi(parts[0])
       poly := make([]Point, k)
       for i := 0; i < k; i++ {
           idx, _ := strconv.Atoi(parts[i+1])
           poly[i] = pts[idx-1]
       }
       // bounding box
       minx, maxx := poly[0].x, poly[0].x
       miny, maxy := poly[0].y, poly[0].y
       for _, p := range poly {
           if p.x < minx {
               minx = p.x
           }
           if p.x > maxx {
               maxx = p.x
           }
           if p.y < miny {
               miny = p.y
           }
           if p.y > maxy {
               maxy = p.y
           }
       }
       cnt := 0
       for _, pt := range pts {
           if pt.x < minx || pt.x > maxx || pt.y < miny || pt.y > maxy {
               continue
           }
           if pointInPoly(poly, pt) {
               cnt++
           }
       }
       out[qi] = cnt
   }
   // print results
   w := bufio.NewWriter(os.Stdout)
   for i, v := range out {
       if i > 0 {
           w.WriteByte(' ')
       }
       fmt.Fprint(w, v)
   }
   w.WriteByte('\n')
   w.Flush()
}
