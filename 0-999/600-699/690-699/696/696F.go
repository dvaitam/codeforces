package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

type R = float64

// com represents a complex number using real x and imaginary y
type com struct{ x, y R }

func (c com) add(d com) com { return com{c.x + d.x, c.y + d.y} }
func (c com) sub(d com) com { return com{c.x - d.x, c.y - d.y} }
func (c com) mulf(f R) com { return com{c.x * f, c.y * f} }
func (c com) divf(f R) com { return com{c.x / f, c.y / f} }
func (c com) abs() R    { return math.Hypot(c.x, c.y) }
// rot90 rotates by 90 degrees (multiply by i)
func (c com) rot90() com { return com{-c.y, c.x} }

func dot(a, b com) R { return a.x*b.x + a.y*b.y }
func det(a, b com) R { return a.x*b.y - a.y*b.x }

const eps = 1e-9

// line in normal form: n·p = c, where n is unit normal
type line struct {
   n com
   c R
}

// newLinePoints constructs a line through p1 and p2
func newLinePoints(p1, p2 com) line {
   d := p2.sub(p1)
   nd := d.divf(d.abs())
   n := nd.rot90()
   return line{n: n, c: dot(p1, n)}
}

// dir returns a direction vector of the line (unit)
func (l line) dir() com { return l.n.rot90() }

// dist returns distance from point p to line l
func dist(l line, p com) R { return math.Abs(dot(p, l.n) - l.c) }

// check_r tests if radius suffices, and sets candB and candL
func check_r(radius R, polygon []com, n int, candB, candL *com) bool {
   i1, m := 2, 0
   best := make([]int, n)
   cand := make([]com, n)
   for i := 0; i < n; i++ {
       cand[i] = polygon[i+1]
   }
   for i := 0; i < n; i++ {
       if i1 < i+2 {
           i1 = i + 2
       }
       if m < i {
           m = i
       }
       best[i] = i1 - i
       l1 := newLinePoints(polygon[i], polygon[i+1])
       for i1 != i+n {
           l2 := newLinePoints(polygon[i1], polygon[i1+1])
           if det(polygon[i+1].sub(polygon[i]), polygon[i1+1].sub(polygon[i1])) < eps {
               break
           }
           for dist(l1, polygon[m+1]) < dist(l2, polygon[m+1]) {
               m++
           }
           low, high := polygon[m], polygon[m+1]
           for k := 0; k < 42; k++ {
               mid := low.add(high).mulf(0.5)
               if dist(l1, mid) < dist(l2, mid) {
                   low = mid
               } else {
                   high = mid
               }
           }
           if dist(l1, low) < radius {
               i1++
               best[i] = i1 - i
               cand[i] = low
           } else {
               break
           }
       }
   }
   for i := 0; i < n; i++ {
       next := (i + best[i]) % n
       if best[i]+best[next] >= n {
           *candB = cand[i]
           *candL = cand[next]
           return true
       }
   }
   return false
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   polygon := make([]com, 2*n)
   for i := 0; i < n; i++ {
       var x, y R
       fmt.Fscan(reader, &x, &y)
       polygon[i] = com{x, y}
   }
   for i := n; i < 2*n; i++ {
       polygon[i] = polygon[i-n]
   }
   start, fin := R(0), R(100000)
   var Bar, Lya com
   for it := 0; it < 84; it++ {
       m := (start + fin) / 2
       var candB, candL com
       if check_r(m, polygon, n, &candB, &candL) {
           fin = m
           Bar = candB
           Lya = candL
       } else {
           start = m
       }
   }
   fmt.Fprintf(writer, "%.10f\n", fin)
   fmt.Fprintf(writer, "%.10f %.10f\n", Bar.x, Bar.y)
   fmt.Fprintf(writer, "%.10f %.10f\n", Lya.x, Lya.y)
}
