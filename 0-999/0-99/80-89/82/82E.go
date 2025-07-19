package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

// pd represents a point or vector with X, Y coordinates
type pd struct { X, Y float64 }

// line represents a line in the form A*x + B*y + C = 0
type line struct { A, B, C float64 }

const EPS = 1e-7

// getLine returns the line passing through points a and b
func getLine(a, b pd) line {
   A := b.Y - a.Y
   B := a.X - b.X
   C := -A*a.X - B*a.Y
   return line{A, B, C}
}

// cp returns the cross product of vectors a and b
func cp(a, b pd) float64 {
   return a.X*b.Y - a.Y*b.X
}

// intersect computes intersection of lines a and b, returns point and true if exists
func intersect(a, b line) (pd, bool) {
   d := cp(pd{a.A, a.B}, pd{b.A, b.B})
   if math.Abs(d) < EPS {
       return pd{}, false
   }
   x := -cp(pd{a.C, a.B}, pd{b.C, b.B}) / d
   y := -cp(pd{a.A, a.C}, pd{b.A, b.C}) / d
   return pd{x, y}, true
}

// dist returns the Euclidean distance between a and b
func dist(a, b pd) float64 {
   dx := a.X - b.X
   dy := a.Y - b.Y
   return math.Sqrt(dx*dx + dy*dy)
}

// equals returns true if points a and b are very close
func equals(a, b pd) bool {
   return dist(a, b) < EPS
}

// getArea computes the area of the polygon given by points p
func getArea(p []pd) float64 {
   if len(p) < 3 {
       return 0
   }
   // copy points
   pts := make([]pd, len(p))
   copy(pts, p)
   // compute centroid c
   var c pd
   n := float64(len(pts))
   for _, v := range pts {
       c.X += v.X / n
       c.Y += v.Y / n
   }
   // shift to centroid
   for i := range pts {
       pts[i].X -= c.X
       pts[i].Y -= c.Y
   }
   // sort by angle
   type ap struct { ang float64; pt pd }
   arr := make([]ap, len(pts))
   for i, v := range pts {
       arr[i] = ap{math.Atan2(v.Y, v.X), v}
   }
   sort.Slice(arr, func(i, j int) bool {
       return arr[i].ang < arr[j].ang
   })
   pts = pts[:0]
   for _, v := range arr {
       pts = append(pts, v.pt)
   }
   // remove duplicates
   uniq := make([]pd, 0, len(pts))
   for _, v := range pts {
       if len(uniq) == 0 || !equals(uniq[len(uniq)-1], v) {
           uniq = append(uniq, v)
       }
   }
   pts = uniq
   // compute area
   res := 0.0
   for i := range pts {
       j := (i + 1) % len(pts)
       ax := pts[i].X - c.X
       ay := pts[i].Y - c.Y
       bx := pts[j].X - c.X
       by := pts[j].Y - c.Y
       res += ax*by - ay*bx
   }
   return math.Abs(res) / 2.0
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var h, f float64
   fmt.Fscan(in, &n, &h, &f)
   // read segments
   s := make([]pd, 0, n*2)
   for i := 0; i < n; i++ {
       var l, r float64
       fmt.Fscan(in, &l, &r)
       if l < 0 && r > 0 {
           s = append(s, pd{l, 0})
           s = append(s, pd{0, r})
       } else {
           s = append(s, pd{l, r})
       }
   }
   cu := pd{0, f}
   cd := pd{0, -f}
   ul := getLine(pd{0, h}, pd{1, h})
   dl := getLine(pd{0, -h}, pd{1, -h})
   ans := 0.0
   // first part
   for i := range s {
       lp := pd{s[i].X, -h}
       rp := pd{s[i].Y, -h}
       lLine := getLine(cd, lp)
       rLine := getLine(cd, rp)
       ulp, _ := intersect(lLine, ul)
       urp, _ := intersect(rLine, ul)
       ans += h * (s[i].Y - s[i].X + urp.X - ulp.X)
   }
   ans *= 2.0
   // overlap corrections
   for i := range s {
       for j := 0; j <= i; j++ {
           pts := make([]pd, 0, 16)
           lp1 := pd{s[i].X, -h}
           rp1 := pd{s[i].Y, -h}
           l1 := getLine(cd, lp1)
           r1 := getLine(cd, rp1)
           ulp1, _ := intersect(l1, ul)
           urp1, _ := intersect(r1, ul)
           lp2 := pd{s[j].X, h}
           rp2 := pd{s[j].Y, h}
           l2 := getLine(cu, lp2)
           r2 := getLine(cu, rp2)
           dlp2, _ := intersect(l2, dl)
           drp2, _ := intersect(r2, dl)
           if ulp1.X >= lp2.X && ulp1.X <= rp2.X {
               pts = append(pts, ulp1)
           }
           if urp1.X >= lp2.X && urp1.X <= rp2.X {
               pts = append(pts, urp1)
           }
           if dlp2.X >= lp1.X && dlp2.X <= rp1.X {
               pts = append(pts, dlp2)
           }
           if drp2.X >= lp1.X && drp2.X <= rp1.X {
               pts = append(pts, drp2)
           }
           if lp2.X >= ulp1.X && lp2.X <= urp1.X {
               pts = append(pts, lp2)
           }
           if rp2.X >= ulp1.X && rp2.X <= urp1.X {
               pts = append(pts, rp2)
           }
           if lp1.X >= dlp2.X && lp1.X <= drp2.X {
               pts = append(pts, lp1)
           }
           if rp1.X >= dlp2.X && rp1.X <= drp2.X {
               pts = append(pts, rp1)
           }
           if cll, ok := intersect(l1, l2); ok && cll.Y >= -h && cll.Y <= h {
               pts = append(pts, cll)
           }
           if clr, ok := intersect(l1, r2); ok && clr.Y >= -h && clr.Y <= h {
               pts = append(pts, clr)
           }
           if crl, ok := intersect(r1, l2); ok && crl.Y >= -h && crl.Y <= h {
               pts = append(pts, crl)
           }
           if crr, ok := intersect(r1, r2); ok && crr.Y >= -h && crr.Y <= h {
               pts = append(pts, crr)
           }
           if i == j {
               ans -= getArea(pts)
           } else {
               ans -= 2.0 * getArea(pts)
           }
       }
   }
   fmt.Printf("%.15f\n", ans)
}
