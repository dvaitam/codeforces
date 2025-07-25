package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Point struct {
   x, y int64
}

// cross product of OA x OB
func cross(o, a, b Point) int64 {
   return (a.x-o.x)*(b.y-o.y) - (a.y-o.y)*(b.x-o.x)
}

// convex hull (monotonic chain), returns points in CCW order, no duplicate of first point at end
func convexHull(a []Point) []Point {
   n := len(a)
   if n <= 1 {
       return append([]Point{}, a...)
   }
   sort.Slice(a, func(i, j int) bool {
       if a[i].x != a[j].x {
           return a[i].x < a[j].x
       }
       return a[i].y < a[j].y
   })
   lower := make([]Point, 0, n)
   for _, p := range a {
       for len(lower) >= 2 && cross(lower[len(lower)-2], lower[len(lower)-1], p) <= 0 {
           lower = lower[:len(lower)-1]
       }
       lower = append(lower, p)
   }
   upper := make([]Point, 0, n)
   for i := n - 1; i >= 0; i-- {
       p := a[i]
       for len(upper) >= 2 && cross(upper[len(upper)-2], upper[len(upper)-1], p) <= 0 {
           upper = upper[:len(upper)-1]
       }
       upper = append(upper, p)
   }
   // drop last point of each (it's the starting point of the other)
   lower = lower[:len(lower)-1]
   upper = upper[:len(upper)-1]
   hull := append(lower, upper...)
   return hull
}

// triple signature for edges
type Sig struct {
   length, cross, dot int64
}

func makeSigs(hull []Point) []Sig {
   h := len(hull)
   if h == 1 {
       return []Sig{{0, 0, 0}}
   }
   // compute edge vectors
   vx := make([]int64, h)
   vy := make([]int64, h)
   for i := 0; i < h; i++ {
       j := (i + 1) % h
       vx[i] = hull[j].x - hull[i].x
       vy[i] = hull[j].y - hull[i].y
   }
   sigs := make([]Sig, h)
   for i := 0; i < h; i++ {
       // length squared
       l := vx[i]*vx[i] + vy[i]*vy[i]
       // angle between edge i and i+1
       ni := (i + 1) % h
       cr := vx[i]*vy[ni] - vy[i]*vx[ni]
       dt := vx[i]*vx[ni] + vy[i]*vy[ni]
       sigs[i] = Sig{l, cr, dt}
   }
   return sigs
}

// KMP prefix function
func prefixFunc(p []Sig) []int {
   n := len(p)
   pi := make([]int, n)
   j := 0
   for i := 1; i < n; i++ {
       for j > 0 && p[i] != p[j] {
           j = pi[j-1]
       }
       if p[i] == p[j] {
           j++
       }
       pi[i] = j
   }
   return pi
}

// check if p is subarray of t via KMP
func kmpSearch(t, p []Sig) bool {
   if len(p) == 0 {
       return true
   }
   pi := prefixFunc(p)
   j := 0
   for i := 0; i < len(t); i++ {
       for j > 0 && t[i] != p[j] {
           j = pi[j-1]
       }
       if t[i] == p[j] {
           j++
           if j == len(p) {
               return true
           }
       }
   }
   return false
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(in, &n, &m)
   A := make([]Point, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &A[i].x, &A[i].y)
   }
   B := make([]Point, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &B[i].x, &B[i].y)
   }
   ha := convexHull(A)
   hb := convexHull(B)
   if len(ha) != len(hb) {
       fmt.Println("NO")
       return
   }
   h := len(ha)
   // special h == 2: segment case or h==1
   if h < 3 {
       // compare segment lengths
       var da, db int64
       da = (ha[1].x-ha[0].x)*(ha[1].x-ha[0].x) + (ha[1].y-ha[0].y)*(ha[1].y-ha[0].y)
       db = (hb[1].x-hb[0].x)*(hb[1].x-hb[0].x) + (hb[1].y-hb[0].y)*(hb[1].y-hb[0].y)
       if da == db {
           fmt.Println("YES")
       } else {
           fmt.Println("NO")
       }
       return
   }
   sa := makeSigs(ha)
   sb := makeSigs(hb)
   // double sb
   t := make([]Sig, 0, len(sb)*2)
   t = append(t, sb...)
   t = append(t, sb...)
   if kmpSearch(t, sa) {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
