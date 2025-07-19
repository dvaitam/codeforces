package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

const N = 100000

type pt struct{ x, y int }

var v []pt

func add(x, y int) { v = append(v, pt{x, y}) }

func vect(a, b, c pt) int64 {
   return int64(b.x-a.x)*int64(c.y-b.y) - int64(b.y-a.y)*int64(c.x-b.x)
}

func dist(a, b pt) float64 {
   dx := float64(a.x - b.x)
   dy := float64(a.y - b.y)
   return math.Hypot(dx, dy)
}

func getr(a, b, c pt) float64 {
   A := dist(a, b)
   B := dist(a, c)
   C := dist(b, c)
   p := (A + B + C) / 2
   area := math.Sqrt(p * (p - A) * (p - B) * (p - C))
   S := 4 * area
   return A * B * C / S
}

func main() {
   in := bufio.NewReader(os.Stdin)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()

   var n int
   fmt.Fscan(in, &n)
   for i := 0; i < n; i++ {
       var x, y, d int
       fmt.Fscan(in, &x, &y, &d)
       if x >= d {
           add(x-d, y)
       } else {
           h := d - x
           if y+h > N {
               add(0, N)
           } else {
               add(0, y+h)
           }
           if y-h < 0 {
               add(0, 0)
           } else {
               add(0, y-h)
           }
       }
       if x+d <= N {
           add(x+d, y)
       } else {
           h := x + d - N
           if y+h > N {
               add(N, N)
           } else {
               add(N, y+h)
           }
           if y-h < 0 {
               add(N, 0)
           } else {
               add(N, y-h)
           }
       }
       if y >= d {
           add(x, y-d)
       } else {
           h := d - y
           if x+h > N {
               add(N, 0)
           } else {
               add(x+h, 0)
           }
           if x-h < 0 {
               add(0, 0)
           } else {
               add(x-h, 0)
           }
       }
       if y+d <= N {
           add(x, y+d)
       } else {
           h := y + d - N
           if x+h > N {
               add(N, N)
           } else {
               add(x+h, N)
           }
           if x-h < 0 {
               add(0, N)
           } else {
               add(0, y-h)
           }
       }
   }
   // sort and unique
   sort.Slice(v, func(i, j int) bool {
       if v[i].x != v[j].x {
           return v[i].x < v[j].x
       }
       return v[i].y < v[j].y
   })
   var u []pt
   for i, p := range v {
       if i == 0 || p.x != v[i-1].x || p.y != v[i-1].y {
           u = append(u, p)
       }
   }
   v = u
   // build convex hull
   var lower, upper []pt
   for _, p := range v {
       for len(lower) > 1 && vect(lower[len(lower)-2], lower[len(lower)-1], p) >= 0 {
           lower = lower[:len(lower)-1]
       }
       for len(upper) > 1 && vect(upper[len(upper)-2], upper[len(upper)-1], p) <= 0 {
           upper = upper[:len(upper)-1]
       }
       lower = append(lower, p)
       upper = append(upper, p)
   }
   // merge hulls
   upper = upper[:len(upper)-1]
   for i, j := 0, len(upper)-1; i < j; i, j = i+1, j-1 {
       upper[i], upper[j] = upper[j], upper[i]
   }
   upper = upper[:len(upper)-1]
   hull := append(lower, upper...)
   // find max circumradius
   bestR := -1e100
   bestIdx := 0
   if len(hull) >= 2 {
       hull = append(hull, hull[0], hull[1])
   }
   for i := 0; i+2 < len(hull); i++ {
       r := getr(hull[i], hull[i+1], hull[i+2])
       if r > bestR {
           bestR = r
           bestIdx = i
       }
   }
   a := hull[bestIdx]
   b := hull[bestIdx+1]
   c := hull[bestIdx+2]
   fmt.Fprintf(w, "%d %d\n", a.x, a.y)
   fmt.Fprintf(w, "%d %d\n", b.x, b.y)
   fmt.Fprintf(w, "%d %d\n", c.x, c.y)
