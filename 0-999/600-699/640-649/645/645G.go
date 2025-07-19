package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

type Point struct {
   x, y float64
}

func (a Point) sub(b Point) Point { return Point{a.x - b.x, a.y - b.y} }
func (a Point) add(b Point) Point { return Point{a.x + b.x, a.y + b.y} }
func (a Point) mul(s float64) Point { return Point{a.x * s, a.y * s} }
func (a Point) div(s float64) Point { return Point{a.x / s, a.y / s} }
func length(a Point) float64 { return math.Hypot(a.x, a.y) }

const eps = 1e-8

// run computes intersection of two circles: center a radius r1, center b radius r2
// returns two intersection points in c1, c2 and ok flag
func run(a, b Point, r1, r2 float64) (c1, c2 Point, ok bool) {
   // ensure r1 <= r2
   if r1 > r2 {
       a, b = b, a
       r1, r2 = r2, r1
   }
   d := length(a.sub(b))
   if r1+r2 < d+eps || d+r1 < r2+eps {
       return Point{}, Point{}, false
   }
   // distance from a to line of intersection
   t := (r1*r1 - r2*r2 + d*d) / (2 * d)
   dir := b.sub(a).div(d)
   mid := a.add(dir.mul(t))
   h2 := r1*r1 - t*t
   if h2 < 0 {
       h2 = 0
   }
   h := math.Sqrt(h2)
   // perpendicular vector
   v := a.sub(b)
   v = Point{-v.y, v.x}.div(length(v))
   c1 = mid.add(v.mul(h))
   c2 = mid.sub(v.mul(h))
   return c1, c2, true
}

type event struct {
   angle float64
   idx   int
}

func ok(points []Point, l0 float64, mid float64) bool {
   n := len(points)
   events := make([]event, 0, 2*n)
   vis := make([]bool, n)
   stack := make([]int, 0, n)
   // center B is (l0,0)
   B := Point{l0, 0}
   C := Point{-l0, 0}
   for i, a := range points {
       // r1: dist from a to C, r2 = mid
       r1 := length(a.sub(C))
       r2 := mid
       c1, c2, ok := run(a, B, r1, r2)
       if !ok {
           continue
       }
       ang1 := math.Atan2(c1.y, c1.x-l0)
       ang2 := math.Atan2(c2.y, c2.x-l0)
       events = append(events, event{ang1, i})
       events = append(events, event{ang2, i})
   }
   if len(events) == 0 {
       return false
   }
   sort.Slice(events, func(i, j int) bool {
       return events[i].angle < events[j].angle
   })
   for _, e := range events {
       if !vis[e.idx] {
           vis[e.idx] = true
           stack = append(stack, e.idx)
       } else {
           top := stack[len(stack)-1]
           if e.idx != top {
               return true
           }
           stack = stack[:len(stack)-1]
       }
   }
   return false
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var l0i int
   if _, err := fmt.Fscan(in, &n, &l0i); err != nil {
       return
   }
   l0 := float64(l0i)
   points := make([]Point, n)
   for i := 0; i < n; i++ {
       var xi, yi int
       fmt.Fscan(in, &xi, &yi)
       points[i] = Point{float64(xi), float64(yi)}
   }
   lo, hi := 0.0, 2*l0
   for it := 0; it < 80; it++ {
       mid := (lo + hi) / 2
       if ok(points, l0, mid) {
           hi = mid
       } else {
           lo = mid
       }
   }
   fmt.Printf("%.12f\n", lo)
}
