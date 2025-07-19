package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

// Point
type P struct { x, y float64 }
func (a P) Add(b P) P  { return P{a.x + b.x, a.y + b.y} }
func (a P) Sub(b P) P  { return P{a.x - b.x, a.y - b.y} }
func (a P) Mul(s float64) P { return P{a.x * s, a.y * s} }
func (a P) Div(s float64) P { return P{a.x / s, a.y / s} }
// cross product (a x b)
func (a P) Cross(b P) float64 { return a.x*b.y - a.y*b.x }

// line through one.two and two.two? Actually holds two points
type L2 struct{ one, two P }
// Intersection of infinite lines
func (l L2) Inter(o L2) P {
   // l.one + dir * t intersects o.one + odir * u
   dir := l.two.Sub(l.one)
   odir := o.two.Sub(o.one)
   den := dir.Cross(odir)
   // assume non-parallel
   t := o.one.Sub(l.one).Cross(odir) / den
   return l.one.Add(dir.Mul(t))
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, q int
   fmt.Fscan(in, &n, &q)
   poly := make([]P, n)
   for i := 0; i < n; i++ {
       var xi, yi float64
       fmt.Fscan(in, &xi, &yi)
       poly[i] = P{xi, yi}
   }
   // compute signed area *2
   var area2 float64
   for i := 0; i < n; i++ {
       j := (i + 1) % n
       area2 += poly[i].Cross(poly[j])
   }
   if area2 > 0 {
       // reverse to make clockwise
       for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
           poly[i], poly[j] = poly[j], poly[i]
       }
       area2 = -area2
   }
   totalArea := math.Abs(area2) / 2
   // build left and right chains
   lowest := 0
   for i := 1; i < n; i++ {
       if poly[i].y < poly[lowest].y {
           lowest = i
       }
   }
   left := []P{}
   for i := lowest; ; {
       left = append(left, poly[i])
       j := (i + 1) % n
       if poly[j].y < poly[i].y {
           break
       }
       i = j
   }
   right := []P{}
   for i := lowest; ; {
       right = append(right, poly[i])
       j := (i + n - 1) % n
       if poly[j].y < poly[i].y {
           break
       }
       i = j
   }
   // build triples
   INF := 1e6 + 5
   triples := make([][3]P, 0)
   triples = append(triples, [3]P{left[0], left[0], left[0]})
   j := 1
   L := len(left)
   R := len(right)
   for i := 1; i < L; i++ {
       for j < R-1 && right[j].y < left[i].y {
           one := L2{left[i], left[i-1]}
           two := L2{right[j], right[j].Sub(P{INF, 0})}
           inter := one.Inter(two)
           mid := inter.Add(right[j]).Div(2)
           triples = append(triples, [3]P{inter, mid, right[j]})
           j++
       }
       if i == L-1 {
           break
       }
       one := L2{right[j], right[j-1]}
       two := L2{left[i], left[i].Add(P{INF, 0})}
       inter := one.Inter(two)
       mid := left[i].Add(inter).Div(2)
       triples = append(triples, [3]P{left[i], mid, inter})
   }
   triples = append(triples, [3]P{left[L-1], left[L-1], left[L-1]})
   // events
   type E struct{ x float64; i, j int }
   events := make([]E, 0, len(triples)*3+q)
   for idx, tri := range triples {
       for j2 := 0; j2 < 3; j2++ {
           events = append(events, E{tri[j2].x, idx, j2})
       }
   }
   for qi := 0; qi < q; qi++ {
       var xv float64
       fmt.Fscan(in, &xv)
       events = append(events, E{xv, -1, qi})
   }
   sort.Slice(events, func(a, b int) bool {
       if events[a].x != events[b].x {
           return events[a].x < events[b].x
       }
       return events[a].i < events[b].i
   })
   on := make([][3]bool, len(triples))
   answer := make([]float64, q)
   var a1, b1, c1 float64
   for _, ev := range events {
       x := ev.x
       if ev.i == -1 {
           s := a1*x*x + b1*x + c1
           answer[ev.j] = s
       } else {
           ti, tj := ev.i, ev.j
           mul := 1.0
           if tj == 1 {
               mul = -2
           } else if tj == 2 {
               mul = 1
           }
           for _, ni := range []int{ti - 1, ti + 1} {
               if ni < 0 || ni >= len(triples) {
                   continue
               }
               A := triples[ti][tj]
               B := triples[ni][tj]
               dx := math.Abs(A.x - B.x)
               if dx >= 1e-7 {
                   a := math.Abs(A.y - B.y) / dx
                   if !on[ni][tj] {
                       tmp := a * mul / 2
                       a1 += tmp
                       b1 -= 2 * tmp * x
                       c1 += tmp * x * x
                   } else {
                       tmp := a * -mul / 2
                       a1 += tmp
                       b1 -= 2 * tmp * B.x
                       c1 += tmp * B.x * B.x
                   }
               }
               if on[ni][tj] {
                   dy := math.Abs(A.y - B.y)
                   dx2 := math.Abs(A.x - B.x)
                   c1 += mul * dy * dx2 / 2
                   b1 += mul * dy
                   c1 -= x * mul * dy
               }
           }
           on[ti][tj] = true
       }
   }
   // output
   for _, v := range answer {
       fmt.Fprintf(out, "%.10f\n", totalArea-v)
   }
}
