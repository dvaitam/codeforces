package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Point represents a 2D integer point
type Point struct {
   x, y int64
}

// Line represents line in form Ax + By + C = 0
type Line struct {
   A, B, C int64
}

func abs64(a int64) int64 {
   if a < 0 {
       return -a
   }
   return a
}

// gcd returns greatest common divisor of a and b
func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

// isIn checks if x between l and r inclusive
func isIn(l, r, x int64) bool {
   if l > r {
       l, r = r, l
   }
   return l <= x && x <= r
}

// isInPoint checks if point (x,y) lies within rectangle defined by a and b
func isInPoint(a, b Point, x, y int64) bool {
   return isIn(a.x, b.x, x) && isIn(a.y, b.y, y)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   // cntMap maps number of pairs t = C(k,2) to k
   cntMap := make(map[int]int)
   for i := 1; i <= n; i++ {
       t := i * (i - 1) / 2
       cntMap[t] = i
   }

   a := make([]Point, n)
   b := make([]Point, n)
   lines := make([]Line, n)
   var ans int64
   for i := 0; i < n; i++ {
       var x1, y1, x2, y2 int64
       fmt.Fscan(in, &x1, &y1, &x2, &y2)
       a[i] = Point{x1, y1}
       b[i] = Point{x2, y2}
       dx := x2 - x1
       dy := y2 - y1
       ans += 1 + gcd(abs64(dx), abs64(dy))
       // line through a[i] and b[i]
       A := y1 - y2
       B := x2 - x1
       C := -A*x1 - B*y1
       lines[i] = Line{A, B, C}
   }
   // collect intersection points (integer and on both segments)
   pts := make([]Point, 0)
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           li, lj := lines[i], lines[j]
           // solve li and lj
           q := lj.A*li.C - lj.C*li.A
           w := lj.C*li.B - lj.B*li.C
           e := li.A*lj.B - lj.A*li.B
           if e != 0 && q%e == 0 && w%e == 0 {
               x := w / e
               y := q / e
               if isInPoint(a[i], b[i], x, y) && isInPoint(a[j], b[j], x, y) {
                   pts = append(pts, Point{x, y})
               }
           }
       }
   }
   // sort points
   sort.Slice(pts, func(i, j int) bool {
       if pts[i].x != pts[j].x {
           return pts[i].x < pts[j].x
       }
       return pts[i].y < pts[j].y
   })
   // subtract overcounts
   m := len(pts)
   for i := 0; i < m; {
       j := i + 1
       for j < m && pts[j] == pts[i] {
           j++
       }
       t := j - i // number of pairs C(k,2)
       if k, ok := cntMap[t]; ok {
           ans -= int64(k - 1)
       }
       i = j
   }
   fmt.Println(ans)
}
