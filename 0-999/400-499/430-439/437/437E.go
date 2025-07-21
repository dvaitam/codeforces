package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

type Point struct {
   x, y int64
}

// orient returns cross product (b-a)x(c-a)
func orient(a, b, c Point) int64 {
   return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}

// intersect checks if segments ab and cd intersect (properly or improperly)
func intersect(a, b, c, d Point) bool {
   // general segment intersection
   o1 := orient(a, b, c)
   o2 := orient(a, b, d)
   o3 := orient(c, d, a)
   o4 := orient(c, d, b)
   if o1 == 0 && onSegment(a, b, c) { return true }
   if o2 == 0 && onSegment(a, b, d) { return true }
   if o3 == 0 && onSegment(c, d, a) { return true }
   if o4 == 0 && onSegment(c, d, b) { return true }
   return (o1 > 0 && o2 < 0 || o1 < 0 && o2 > 0) &&
          (o3 > 0 && o4 < 0 || o3 < 0 && o4 > 0)
}

// onSegment checks if c lies on segment ab
func onSegment(a, b, c Point) bool {
   return c.x >= min(a.x, b.x) && c.x <= max(a.x, b.x) &&
          c.y >= min(a.y, b.y) && c.y <= max(a.y, b.y)
}

func min(a, b int64) int64 { if a < b { return a }; return b }
func max(a, b int64) int64 { if a > b { return a }; return b }

// pointInPoly returns true if point p is strictly inside polygon poly
func pointInPoly(poly []Point, p struct{ x, y float64 }) bool {
   cnt := 0
   n := len(poly)
   for i := 0; i < n; i++ {
       a := poly[i]
       b := poly[(i+1)%n]
       ay := float64(a.y)
       by := float64(b.y)
       if ay > by {
           ay, by = by, ay
           // swap a,b roles for x calculation
           // but we only need y ordering here
       }
       if p.y <= ay || p.y > by {
           continue
       }
       // compute intersection x coordinate
       ax := float64(a.x)
       bx := float64(b.x)
       xint := ax + (bx-ax)*( (p.y - ay) / (by - ay) )
       if xint > p.x {
           cnt++
       }
   }
   return cnt%2 == 1
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   poly := make([]Point, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &poly[i].x, &poly[i].y)
   }
   // compute polygon orientation
   var area2 int64
   for i := 0; i < n; i++ {
       j := (i + 1) % n
       area2 += poly[i].x*poly[j].y - poly[i].y*poly[j].x
   }
   ccw := area2 > 0
   // valid[i][j] diagonal or edge
   valid := make([][]bool, n)
   for i := range valid {
       valid[i] = make([]bool, n)
   }
   // adjacency
   for i := 0; i < n; i++ {
       valid[i][(i+1)%n] = true
       valid[(i+1)%n][i] = true
   }
   // check all possible diagonals i<j
   for i := 0; i < n; i++ {
       for j := i + 2; j < n; j++ {
           if i == 0 && j == n-1 {
               valid[i][j] = true
               valid[j][i] = true
               continue
           }
           a := poly[i]
           b := poly[j]
           ok := true
           // check intersections
           for k := 0; k < n; k++ {
               k2 := (k + 1) % n
               // skip edges incident to i or j
               if k == i || k == j || k2 == i || k2 == j {
                   continue
               }
               if intersect(a, b, poly[k], poly[k2]) {
                   ok = false
                   break
               }
           }
           if !ok {
               continue
           }
           // check midpoint inside
           mx := (float64(a.x) + float64(b.x)) * 0.5
           my := (float64(a.y) + float64(b.y)) * 0.5
           if !pointInPoly(poly, struct{ x, y float64 }{mx, my}) {
               continue
           }
           valid[i][j] = true
           valid[j][i] = true
       }
   }
   // dp[i][j]
   dp := make([][]int64, n)
   for i := range dp {
       dp[i] = make([]int64, n)
       dp[i][i] = 1
       if i+1 < n {
           dp[i][i+1] = 1
       }
   }
   for length := 2; length < n; length++ {
       for i := 0; i+length < n; i++ {
           j := i + length
           if !valid[i][j] {
               dp[i][j] = 0
               continue
           }
           var ways int64
           for k := i + 1; k < j; k++ {
               if !valid[i][k] || !valid[k][j] {
                   continue
               }
               o := orient(poly[i], poly[k], poly[j])
               if (ccw && o <= 0) || (!ccw && o >= 0) {
                   continue
               }
               ways = (ways + dp[i][k]*dp[k][j]) % mod
           }
           dp[i][j] = ways
       }
   }
   fmt.Println(dp[0][n-1])
}
