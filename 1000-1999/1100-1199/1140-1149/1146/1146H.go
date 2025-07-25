package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Point represents a 2D point
type Point struct { x, y int64 }

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   pts := make([]Point, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &pts[i].x, &pts[i].y)
   }
   // sort by y, then x
   sort.Slice(pts, func(i, j int) bool {
       if pts[i].y != pts[j].y {
           return pts[i].y < pts[j].y
       }
       return pts[i].x < pts[j].x
   })
   var ans int64
   // temporary arrays
   for s := 0; s < n; s++ {
       // pivot pts[s]
       m := n - s - 1
       if m < 4 {
           break
       }
       O := pts[s]
       Q := make([]Point, m)
       for i := 0; i < m; i++ {
           Q[i] = pts[s+1+i]
       }
       // sort Q by angle around O (upper half-plane guaranteed)
       sort.Slice(Q, func(i, j int) bool {
           return (Q[i].x-O.x)*(Q[j].y-O.y) - (Q[i].y-O.y)*(Q[j].x-O.x) > 0
       })
       // dp3[k][j]: number of chains O->Q[i]->Q[k]->Q[j]
       dp3 := make([][]int, m)
       for i := range dp3 {
           dp3[i] = make([]int, m)
       }
       // build dp3
       for k := 1; k < m; k++ {
           for j := k + 1; j < m; j++ {
               cnt := 0
               // count i<k with orient(Q[i],Q[k],Q[j])>0
               xk, yk := Q[k].x, Q[k].y
               xj, yj := Q[j].x, Q[j].y
               for i := 0; i < k; i++ {
                   // cross of (Q[k]-Q[i]) x (Q[j]-Q[i]) positive means CCW at i
                   // but we need orient(Q[i],Q[k],Q[j])>0
                   if (xk-Q[i].x)*(yj-Q[i].y)-(yk-Q[i].y)*(xj-Q[i].x) > 0 {
                       cnt++
                   }
               }
               dp3[k][j] = cnt
           }
       }
       // dp4[u][t]: number of chains of 4 edges O->...->Q[u]->Q[t]
       dp4 := make([][]int, m)
       for i := range dp4 {
           dp4[i] = make([]int, m)
       }
       for u := 2; u < m; u++ {
           xu, yu := Q[u].x, Q[u].y
           for t := u + 1; t < m; t++ {
               xt, yt := Q[t].x, Q[t].y
               cnt := 0
               for k := 1; k < u; k++ {
                   // orient(Q[k],Q[u],Q[t])>0
                   if (xu-Q[k].x)*(yt-Q[k].y)-(yu-Q[k].y)*(xt-Q[k].x) > 0 {
                       cnt += dp3[k][u]
                   }
               }
               dp4[u][t] = cnt
           }
       }
       // closure: sum dp4[u][t] where orient(Q[u],Q[t],O)>0
       for u := 2; u < m; u++ {
           xu, yu := Q[u].x, Q[u].y
           for t := u + 1; t < m; t++ {
               if dp4[u][t] > 0 {
                   // orient(Q[u],Q[t],O)>0
                   if (Q[t].x-xu)*(O.y-yu)-(Q[t].y-yu)*(O.x-xu) > 0 {
                       ans += int64(dp4[u][t])
                   }
               }
           }
       }
   }
   // output result
   fmt.Println(ans)
}
