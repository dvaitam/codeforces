package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

const eps = 1e-9

// Pt represents a point or vector.
type Pt struct { x, y float64 }

func (a Pt) Add(b Pt) Pt { return Pt{a.x + b.x, a.y + b.y} }
func (a Pt) Sub(b Pt) Pt { return Pt{a.x - b.x, a.y - b.y} }
func (a Pt) Mul(f float64) Pt { return Pt{a.x * f, a.y * f} }
func (a Pt) Dot(b Pt) float64 { return a.x*b.x + a.y*b.y }
func (a Pt) Cross(b Pt) float64 { return a.x*b.y - a.y*b.x }
func (a Pt) Abs() float64 { return math.Hypot(a.x, a.y) }
// Up returns true if vector is in upper half-plane or on positive x-axis.
func (a Pt) Up() bool {
   if math.Abs(a.y) < eps {
       return a.x > 0
   }
   return a.y > 0
}

// line represents line: v . p = c, where v is unit normal.
type line struct {
   v Pt
   c float64
}

// newLinePoints from two points p1->p2.
func newLinePoints(p1, p2 Pt) line {
   // normal is perpendicular to direction
   dir := p2.Sub(p1)
   n := Pt{-dir.y, dir.x}
   d := n.Abs()
   if d > 0 {
       n = n.Mul(1.0 / d)
   }
   return line{v: n, c: n.Dot(p1)}
}

// signedDist returns v.p - c
func (l line) signedDist(p Pt) float64 {
   return l.v.Dot(p) - l.c
}

// cmpAngle orders vectors by angle around origin.
func cmpAngle(a, b Pt) bool {
   au, bu := a.Up(), b.Up()
   if au != bu {
       return au
   }
   return a.Cross(b) > eps
}

// eqLine checks if two lines have same direction and coincide.
func eqLine(a, b line) bool {
   if a.v.Up() != b.v.Up() {
       return false
   }
   return math.Abs(a.v.Cross(b.v)) < eps
}

// cmpLine for sorting half-planes
func cmpLine(a, b line) bool {
   au, bu := a.v.Up(), b.v.Up()
   if au != bu {
       return au
   }
   cr := a.v.Cross(b.v)
   if math.Abs(cr) > eps {
       return cr > 0
   }
   return a.c > b.c
}

// det3x3 of lines a,b,c
func det3x3(a, b, c line) float64 {
   return a.c*(b.v.Cross(c.v)) + b.c*(c.v.Cross(a.v)) + c.c*(a.v.Cross(b.v))
}

// intersect two lines
func intersect(l1, l2 line) Pt {
   d := l1.v.x*l2.v.y - l1.v.y*l2.v.x
   if math.Abs(d) < eps {
       return Pt{1e18, 1e18}
   }
   dx := l1.c*l2.v.y - l1.v.y*l2.c
   dy := l1.v.x*l2.c - l1.c*l2.v.x
   return Pt{dx / d, dy / d}
}

// halfplanesIntersection returns polygon of intersection (possibly empty)
func halfplanesIntersection(ls []line) []Pt {
   sort.Slice(ls, func(i, j int) bool { return cmpLine(ls[i], ls[j]) })
   // unique
   uniq := make([]line, 0, len(ls))
   for i, l := range ls {
       if i == 0 || !eqLine(ls[i-1], l) {
           uniq = append(uniq, l)
       }
   }
   ls = uniq
   n := len(ls)
   st := make([]int, 0, n*2)
   // double pass
   for iter := 0; iter < 2; iter++ {
       for i := 0; i < n; i++ {
           for len(st) > 1 {
               j := st[len(st)-1]
               k := st[len(st)-2]
               if ls[k].v.Cross(ls[i].v) <= eps || det3x3(ls[k], ls[j], ls[i]) <= eps {
                   break
               }
               st = st[:len(st)-1]
           }
           st = append(st, i)
       }
   }
   pos := make([]int, n)
   for i := range pos {
       pos[i] = -1
   }
   ok := false
   var seq []int
   for i, id := range st {
       if pos[id] != -1 {
           seq = st[pos[id]:i]
           ok = true
           break
       }
       pos[id] = i
   }
   if !ok {
       return nil
   }
   // build intersection points
   k := len(seq)
   res := make([]Pt, k)
   M := Pt{0, 0}
   for i := 0; i < k; i++ {
       l1 := ls[seq[i]]
       l2 := ls[seq[(i+1)%k]]
       p := intersect(l1, l2)
       res[i] = p
       M = M.Add(p)
   }
   M = M.Mul(1.0 / float64(k))
   for _, id := range seq {
       if ls[id].signedDist(M) < -eps {
           return nil
       }
   }
   return res
}

func solve(in *bufio.Reader, out *bufio.Writer) {
   var tn int
   fmt.Fscan(in, &tn)
   for ; tn > 0; tn-- {
       var n int
       var px, py float64
       fmt.Fscan(in, &n, &px, &py)
       pivot := Pt{px, py}
       // read points
       ev := make([]Pt, 0, (n-1)*2)
       for i := 1; i < n; i++ {
           var x, y float64
           fmt.Fscan(in, &x, &y)
           v := Pt{x, y}.Sub(pivot)
           ev = append(ev, v)
           ev = append(ev, Pt{-v.x, -v.y})
       }
       sort.Slice(ev, func(i, j int) bool { return cmpAngle(ev[i], ev[j]) })
       // bounding box
       const B = 1e6
       hp := make([]line, 0, len(ev)+4)
       negPivot := pivot.Mul(-1)
       box := []Pt{
           negPivot.Add(Pt{B, B}),
           negPivot.Add(Pt{-B, B}),
           negPivot.Add(Pt{-B, -B}),
           negPivot.Add(Pt{B, -B}),
       }
       for i := 0; i < 4; i++ {
           hp = append(hp, newLinePoints(box[i], box[(i+1)%4]))
       }
       m := len(ev)
       bad := false
       for i := 0; i < m; i++ {
           A := ev[i]
           Bv := ev[(i+1)%m]
           if math.Abs(A.Cross(Bv)) < eps {
               fmt.Fprintln(out, 0)
               bad = true
               break
           }
           l := newLinePoints(A, Bv)
           // ensure origin in positive side
           if l.signedDist(Pt{0, 0}) < 0 {
               l = newLinePoints(Bv, A)
           }
           hp = append(hp, l)
       }
       if bad {
           continue
       }
       poly := halfplanesIntersection(hp)
       var area float64
       k := len(poly)
       for i := 0; i < k; i++ {
           area += poly[i].Cross(poly[(i+1)%k])
       }
       area = math.Abs(area) / 2.0
       fmt.Fprintf(out, "%.10f\n", area)
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   solve(in, out)
}
