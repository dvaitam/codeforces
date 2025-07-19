package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Point represents a point with coordinates and original index
type Point struct {
   x, y int64
   idx  int
}

// Line represents a directed line between two points
type Line struct {
   dx, dy     int64
   i1, i2     int
}

var (
   a   []Point
   rev []int
   v   []Line
   n   int
   p2  int64
)

// abs returns the absolute value of x
func abs(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

// ccw returns twice the area of triangle abc (absolute value)
func ccw(a, b, c Point) int64 {
   dx1 := b.x - a.x
   dy1 := b.y - a.y
   dx2 := c.x - a.x
   dy2 := c.y - a.y
   return abs(dx1*dy2 - dy1*dx2)
}

// report outputs the result and exits
func report(x, y, z int) {
   fmt.Println("Yes")
   fmt.Println(a[x].x, a[x].y)
   fmt.Println(a[y].x, a[y].y)
   fmt.Println(a[z].x, a[z].y)
   os.Exit(0)
}

// solve searches for a third point forming area == p2 with edge c1-c2
func solve(c1, c2 int, l int64) {
   // search in [c2..n-1]
   s, e := c2, n-1
   for s != e {
       m := (s + e + 1) / 2
       if ccw(a[c1], a[c2], a[m]) <= l {
           s = m
       } else {
           e = m - 1
       }
   }
   if ccw(a[c1], a[c2], a[s]) == l {
       report(c1, c2, s)
   }
   // search in [0..c1]
   s, e = 0, c1
   for s != e {
       m := (s + e) / 2
       if ccw(a[c1], a[c2], a[m]) <= l {
           e = m
       } else {
           s = m + 1
       }
   }
   if ccw(a[c1], a[c2], a[s]) == l {
       report(c1, c2, s)
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var p int64
   fmt.Fscan(reader, &n, &p)
   p2 = p << 1
   a = make([]Point, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i].x, &a[i].y)
       a[i].idx = i
   }
   // sort points by x, then y
   sort.Slice(a, func(i, j int) bool {
       if a[i].x != a[j].x {
           return a[i].x < a[j].x
       }
       return a[i].y < a[j].y
   })
   // initialize reverse mapping
   rev = make([]int, n)
   for i := 0; i < n; i++ {
       rev[a[i].idx] = i
   }
   // build all directed edges
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           v = append(v, Line{
               dx: a[j].x - a[i].x,
               dy: a[j].y - a[i].y,
               i1: a[i].idx,
               i2: a[j].idx,
           })
       }
   }
   // sort edges by angle, then by indices
   sort.Slice(v, func(i, j int) bool {
       a1 := v[i]
       b1 := v[j]
       cw := a1.dx*b1.dy - b1.dx*a1.dy
       if cw != 0 {
           return cw > 0
       }
       if a1.i1 != b1.i1 {
           return a1.i1 < b1.i1
       }
       return a1.i2 < b1.i2
   })
   // sweep edges
   for _, line := range v {
       c1 := rev[line.i1]
       c2 := rev[line.i2]
       if c1 > c2 {
           c1, c2 = c2, c1
       }
       solve(c1, c2, p2)
       // swap positions in a and rev
       a[c1], a[c2] = a[c2], a[c1]
       rev[line.i1], rev[line.i2] = rev[line.i2], rev[line.i1]
   }
   fmt.Println("No")
}
