package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

type IP struct { X, Y int64 }
type P struct { X, Y float64 }

func detIP(a, b IP) int64 { return a.X*b.Y - a.Y*b.X }
func dotIP(a, b IP) int64 { return a.X*b.X + a.Y*b.Y }
func subIP(a, b IP) IP { return IP{a.X - b.X, a.Y - b.Y} }
func sgn64(x int64) int { if x > 0 { return 1 } else if x < 0 { return -1 } else { return 0 } }

// distance between two P
func distP(a, b P) float64 {
   dx := a.X - b.X
   dy := a.Y - b.Y
   return math.Hypot(dx, dy)
}

// line represented by two points (as P)
type line struct { S, T P }

func hasInter(a, b line) bool {
   // determinant signs differ
   // compute using floats? we can use floats since inputs are ints
   s1 := (a.T.X - a.S.X)*(b.S.Y - a.S.Y) - (a.T.Y - a.S.Y)*(b.S.X - a.S.X)
   s2 := (a.T.X - a.S.X)*(b.T.Y - a.S.Y) - (a.T.Y - a.S.Y)*(b.T.X - a.S.X)
   return s1 != s2
}

func lineInter(a, b line) P {
   // intersection point of lines a and b
   // using determinant form
   // convert to IP-like floats
   bs := b.S
   bt := b.T
   // compute s1 = det(a.T-a.S, b.S-a.S)
   s1 := (a.T.X - a.S.X)*(bs.Y - a.S.Y) - (a.T.Y - a.S.Y)*(bs.X - a.S.X)
   s2 := (a.T.X - a.S.X)*(bt.Y - a.S.Y) - (a.T.Y - a.S.Y)*(bt.X - a.S.X)
   // P = (b.S * s2 - b.T * s1) / (s2 - s1)
   den := s2 - s1
   return P{
       X: (bs.X*s2 - bt.X*s1) / den,
       Y: (bs.Y*s2 - bt.Y*s1) / den,
   }
}

func turnLeft(u, x, y IP) bool {
   return detIP(subIP(x, u), subIP(y, u)) > 0
}

func pointOnSeg(p IP, a, b IP) bool {
   if sgn64(detIP(subIP(p, a), subIP(b, a))) != 0 {
       return false
   }
   return dotIP(subIP(a, p), subIP(b, p)) <= 0
}

// find index maximizing predicate f in convex polygon order
func findMax(pts []IP, f func(IP, IP) bool) int {
   l, r := 0, len(pts)-1
   // compute d = !f(pts[l], pts[r])
   d := 0
   if !f(pts[l], pts[r]) {
       d = 1
   }
   if d == 1 {
       l, r = r, l
   }
   for ;; {
       diff := r - l
       if diff < 0 {
           diff = -diff
       }
       if diff <= 1 {
           break
       }
       mid := (l + r + d) / 2
       var next int
       if d == 1 {
           next = mid + 1
       } else {
           next = mid - 1
       }
       if f(pts[mid], pts[l]) && f(pts[mid], pts[next]) {
           l = mid
       } else {
           r = mid
       }
   }
   return l
}

func getTan(u IP, pts []IP, id map[IP]int) [2]IP {
   n := len(pts)
   if idx, ok := id[u]; ok {
       return [2]IP{pts[(idx-1+n)%n], pts[(idx+1)%n]}
   }
   // check on segment between first and last
   if pointOnSeg(u, pts[0], pts[n-1]) {
       return [2]IP{pts[n-1], pts[0]}
   }
   t1 := findMax(pts, func(x, y IP) bool { return turnLeft(u, y, x) })
   t2 := findMax(pts, func(x, y IP) bool { return turnLeft(u, x, y) })
   return [2]IP{pts[t1], pts[t2]}
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   pts := make([]IP, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &pts[i].X, &pts[i].Y)
   }
   id := make(map[IP]int, n)
   for i, p := range pts {
       id[p] = i
   }
   var Q int
   fmt.Fscan(in, &Q)
   for qi := 0; qi < Q; qi++ {
       var u, v IP
       fmt.Fscan(in, &u.X, &u.Y, &v.X, &v.Y)
       uu := getTan(u, pts, id)
       vv := getTan(v, pts, id)
       // check if both outside
       cond := detIP(subIP(uu[1], u), subIP(v, u)) > 0 && detIP(subIP(v, u), subIP(uu[0], u)) > 0 &&
           detIP(subIP(vv[1], v), subIP(u, v)) > 0 && detIP(subIP(u, v), subIP(vv[0], v)) > 0
       ans := 1e100
       if cond {
           // try intersections
           uP := P{float64(u.X), float64(u.Y)}
           vP := P{float64(v.X), float64(v.Y)}
           for _, pu := range uu {
               for _, pv := range vv {
                   la := line{uP, P{float64(pu.X), float64(pu.Y)}}
                   lb := line{vP, P{float64(pv.X), float64(pv.Y)}}
                   if hasInter(la, lb) {
                       ip := lineInter(la, lb)
                       d := distP(uP, ip) + distP(vP, ip)
                       if d < ans {
                           ans = d
                       }
                   }
               }
           }
       } else {
           // direct
           uP := P{float64(u.X), float64(u.Y)}
           vP := P{float64(v.X), float64(v.Y)}
           ans = distP(uP, vP)
       }
       if ans < 1e100 {
           fmt.Fprintf(out, "%.9f\n", ans)
       } else {
           out.WriteString("-1\n")
       }
   }
}
