package main

import (
   "bufio"
   "fmt"
   "math"
   "math/rand"
   "os"
)

const (
   EPS = 1e-9
   PI  = math.Pi
)

// point or vector
type pt struct {
   x, y float64
}

// Less compares lex order with EPS tolerance
func (a pt) Less(b pt) bool {
   if a.x < b.x-EPS {
       return true
   }
   if math.Abs(a.x-b.x) < EPS && a.y < b.y-EPS {
       return true
   }
   return false
}

// line represented as ax + by + c = 0
type line struct {
   a, b, c float64
}

// construct line through p and q
func lineFrom(p, q pt) line {
   a := p.y - q.y
   b := q.x - p.x
   c := -a*p.x - b*p.y
   l := line{a, b, c}
   l.norm()
   return l
}

// normalize line coefficients
func (l *line) norm() {
   z := math.Hypot(l.a, l.b)
   if math.Abs(z) > EPS {
       l.a /= z
       l.b /= z
       l.c /= z
   }
}

// signed distance from point to line
func (l *line) dist(p pt) float64 {
   return l.a*p.x + l.b*p.y + l.c
}

// determinant 2x2
func det(a, b, c, d float64) float64 {
   return a*d - b*c
}

// cross product of (a) x (b)
func detPt(a, b pt) float64 {
   return det(a.x, a.y, b.x, b.y)
}

// check if x is between l and r
func betw(l, r, x float64) bool {
   return math.Min(l, r) <= x+EPS && x <= math.Max(l, r)+EPS
}

// check 1D overlap
func intersect1d(a, b, c, d float64) bool {
   if a > b {
       a, b = b, a
   }
   if c > d {
       c, d = d, c
   }
   return math.Max(a, c) <= math.Min(b, d)+EPS
}

// intersect segments ab and cd
func intersect(a, b, c, d pt) (bool, pt, pt) {
   if !intersect1d(a.x, b.x, c.x, d.x) || !intersect1d(a.y, b.y, c.y, d.y) {
       return false, pt{}, pt{}
   }
   m := lineFrom(a, b)
   n := lineFrom(c, d)
   zn := det(m.a, m.b, n.a, n.b)
   if math.Abs(zn) < EPS {
       if math.Abs(m.dist(c)) > EPS || math.Abs(n.dist(a)) > EPS {
           return false, pt{}, pt{}
       }
       if b.Less(a) {
           a, b = b, a
       }
       if d.Less(c) {
           c, d = d, c
       }
       // max in lex order
       var left, right pt
       if a.Less(c) {
           left = c
       } else {
           left = a
       }
       if b.Less(d) {
           right = b
       } else {
           right = d
       }
       return true, left, right
   }
   x := -det(m.c, m.b, n.c, n.b) / zn
   y := -det(m.a, m.c, n.a, n.c) / zn
   ip := pt{x, y}
   if betw(a.x, b.x, ip.x) && betw(a.y, b.y, ip.y) && betw(c.x, d.x, ip.x) && betw(c.y, d.y, ip.y) {
       return true, ip, ip
   }
   return false, pt{}, pt{}
}

// rotate vector v by angle
func rotate(v pt, angle float64) pt {
   return pt{
       v.x*math.Cos(angle) - v.y*math.Sin(angle),
       v.x*math.Sin(angle) + v.y*math.Cos(angle),
   }
}

// distance between points
func dist(a, b pt) float64 {
   dx := a.x - b.x
   dy := a.y - b.y
   return math.Hypot(dx, dy)
}

// compute Fermat point of triangle a,b,c
func fermat(a, b, c pt) pt {
   // ensure orientation
   if detPt(pt{b.x - a.x, b.y - a.y}, pt{c.x - a.x, c.y - a.y}) > 0 {
       b, c = c, b
   }
   // helper: build segment [pkt, x]
   f := func(x, y, z pt) (pt, pt) {
       bok := pt{z.x - y.x, z.y - y.y}
       bok = rotate(bok, PI/3.0)
       pkt := pt{y.x + bok.x, y.y + bok.y}
       return pkt, x
   }
   o1p, o1q := f(a, b, c)
   o2p, o2q := f(b, c, a)
   ok, A, _ := intersect(o1p, o1q, o2p, o2q)
   if !ok {
       pts := []pt{a, b, c}
       best := pts[0]
       minLen := math.Inf(1)
       for i := 0; i < 3; i++ {
           p := pts[i]
           prev := pts[(i+2)%3]
           next := pts[(i+1)%3]
           now := dist(p, prev) + dist(p, next)
           if now < minLen {
               minLen = now
               best = p
           }
       }
       return best
   }
   return A
}

// evaluate maximal perimeter over triangles formed by a and each edge
func ans(a pt, PT [3]pt) float64 {
   ret := 0.0
   for i := 0; i < 3; i++ {
       s := fermat(a, PT[i], PT[(i+1)%3])
       val := dist(s, a) + dist(s, PT[i]) + dist(s, PT[(i+1)%3])
       if val > ret {
           ret = val
       }
   }
   return ret
}

func solve() {
   reader := bufio.NewReader(os.Stdin)
   var PT [3]pt
   for i := 0; i < 3; i++ {
       fmt.Fscan(reader, &PT[i].x, &PT[i].y)
   }
   rand.Seed(123)
   // initial point
   akt := fermat(PT[0], PT[1], PT[2])
   best := ans(akt, PT)
   for eps := 1e5; eps >= 1e-9; eps /= 1.01 {
       ang := rand.Float64() * 2 * PI
       nowy := pt{akt.x + eps*math.Cos(ang), akt.y + eps*math.Sin(ang)}
       w := ans(nowy, PT)
       if w < best {
           best = w
           akt = nowy
       }
   }
   fmt.Printf("%.10f\n", best)
}

func main() {
   solve()
}
