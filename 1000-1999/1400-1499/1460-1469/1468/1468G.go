package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

// P represents a point or vector in 2D space
type P struct { x, y float64 }

func sub(a, b P) P { return P{a.x - b.x, a.y - b.y} }
func cross(a, b P) float64 { return a.x*b.y - a.y*b.x }
func dis(a P) float64 { return a.x*a.x + a.y*a.y }

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var H float64
   if _, err := fmt.Fscan(reader, &n, &H); err != nil {
       return
   }
   a := make([]P, n+1)
   for i := 1; i <= n; i++ {
       var u, v float64
       fmt.Fscan(reader, &u, &v)
       a[i] = P{u, v}
   }
   // Origin O is the last point lifted by H in y-direction
   O := a[n]
   O.y += H
   id := n
   ans := 0.0
   // Traverse points in reverse
   for i := n - 1; i >= 1; i-- {
       tmp0 := sub(a[id], O)
       tmp1 := sub(a[i+1], O)
       tmp2 := sub(a[i], O)
       c1 := cross(tmp1, tmp0)
       c2 := cross(tmp2, tmp0)
       if c1 >= 0 && c2 >= 0 {
           // Direct visible segment
           dx := a[i].x - a[i+1].x
           dy := a[i].y - a[i+1].y
           ans += math.Hypot(dx, dy)
       } else if c2 >= 0 {
           // Compute intersection with line of sight
           K := tmp0.y / tmp0.x
           ke := (tmp1.y - tmp2.y) / (tmp1.x - tmp2.x)
           be := tmp1.y - ke*tmp1.x
           rx := be / (K - ke)
           ry := K * rx
           // Map back to original coordinates
           r := P{rx + O.x, ry + O.y}
           now := dis(sub(a[i], r))
           if now > 0 {
               ans += math.Sqrt(now)
           }
       }
       if c2 >= 0 {
           id = i
       }
   }
   fmt.Printf("%.10f\n", ans)
}
