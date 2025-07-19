package main

import (
   "bufio"
   "fmt"
   "os"
   "math"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var x1, y1, x2, y2 float64
   fmt.Fscan(reader, &x1, &y1, &x2, &y2)
   const inf = 1e10
   l, r := 0.0, inf
   for i := 0; i < n; i++ {
       var x, y, vx, vy float64
       fmt.Fscan(reader, &x, &y, &vx, &vy)
       // X-axis interval
       var u, v float64
       if vx == 0 {
           if x > x1 && x < x2 {
               u = 0
               v = inf
           } else {
               fmt.Println(-1)
               return
           }
       } else {
           u = (x1 - x) / vx
           v = (x2 - x) / vx
           if u > v {
               u, v = v, u
           }
       }
       l = math.Max(l, u)
       r = math.Min(r, v)
       // Y-axis interval
       if vy == 0 {
           if y > y1 && y < y2 {
               u = 0
               v = inf
           } else {
               fmt.Println(-1)
               return
           }
       } else {
           u = (y1 - y) / vy
           v = (y2 - y) / vy
           if u > v {
               u, v = v, u
           }
       }
       l = math.Max(l, u)
       r = math.Min(r, v)
   }
   if l >= r {
       fmt.Println(-1)
   } else {
       fmt.Printf("%.10f\n", l)
   }
}
