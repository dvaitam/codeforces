package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

// Point represents a 2D point with original index
type Point struct {
   x, y, ind int
}

// distance returns Euclidean distance between a and b
func distance(a, b Point) float64 {
   dx := float64(a.x - b.x)
   dy := float64(a.y - b.y)
   return math.Hypot(dx, dy)
}

// collinear checks if points a, b, c are collinear
func collinear(a, b, c Point) bool {
   x1 := int64(b.x - a.x)
   y1 := int64(b.y - a.y)
   x2 := int64(c.x - a.x)
   y2 := int64(c.y - a.y)
   return x1*y2 == x2*y1
}

// minFloat returns the minimum of two floats
func minFloat(a, b float64) float64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(reader, &n, &k)
   k--
   v := make([]Point, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &v[i].x, &v[i].y)
       v[i].ind = i
   }
   if n > 3 {
       for i := 0; i < 4; i++ {
           if collinear(v[0], v[1], v[2]) {
               break
           }
           // rotate first 4 points
           temp := v[0]
           copy(v[0:3], v[1:4])
           v[3] = temp
       }
   }
   for i := 3; i < n; i++ {
       if !collinear(v[0], v[1], v[i]) {
           v[0], v[i] = v[i], v[0]
           break
       }
   }
   // sort remaining points by x, then y
   sort.Slice(v[1:], func(i, j int) bool {
       a := v[1+i]
       b := v[1+j]
       if a.x != b.x {
           return a.x < b.x
       }
       return a.y < b.y
   })
   // find target index
   j := 0
   for idx := range v {
       if v[idx].ind == k {
           j = idx
           break
       }
   }
   // handle simple case
   if j == 0 {
       res := minFloat(
           distance(v[0], v[1])+distance(v[1], v[n-1]),
           distance(v[0], v[n-1])+distance(v[1], v[n-1]),
       )
       fmt.Printf("%.10f\n", res)
       return
   }
   const inf = 1e18
   meg := inf
   // candidate 1
   var cand, t1, t2 float64
   t1 = distance(v[1], v[0]) + distance(v[1], v[j])
   if j == n-1 {
       t2 = 0.0
   } else {
       t2 = minFloat(distance(v[0], v[j+1]), distance(v[0], v[n-1]))
   }
   if j == n-1 {
       cand = t1 + t2
   } else {
       cand = t1 + t2 + distance(v[j+1], v[n-1])
   }
   meg = minFloat(meg, cand)
   // candidate 2
   t1 = distance(v[n-1], v[0]) + distance(v[j], v[n-1])
   if j == 1 {
       t2 = 0.0
   } else {
       t2 = minFloat(distance(v[0], v[j-1]), distance(v[0], v[1]))
   }
   if j == 1 {
       cand = t1 + t2
   } else {
       cand = t1 + t2 + distance(v[j-1], v[1])
   }
   meg = minFloat(meg, cand)
   // candidates 3
   for i := 1; i <= j; i++ {
       base := distance(v[j], v[n-1])*2 + distance(v[i], v[j]) + distance(v[i], v[0])
       if i == 1 {
           cand = base
       } else {
           cand = base + minFloat(distance(v[i-1], v[0]), distance(v[1], v[0])) + distance(v[1], v[i-1])
       }
       meg = minFloat(meg, cand)
   }
   // candidates 4
   for i := j; i < n; i++ {
       base := distance(v[1], v[j])*2 + distance(v[j], v[i]) + distance(v[i], v[0])
       if i == n-1 {
           cand = base
       } else {
           cand = base + minFloat(distance(v[i+1], v[0]), distance(v[n-1], v[0])) + distance(v[i+1], v[n-1])
       }
       meg = minFloat(meg, cand)
   }
   fmt.Printf("%.10f\n", meg)
}
