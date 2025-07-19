package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

const INF = 1000000007

func minInt(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func dist2(x1, y1, x2, y2, x3, y3, x4, y4 int) int {
   // squared distance between rectangles [x1,y1]-[x2,y2] and [x3,y3]-[x4,y4]
   dx, dy := 0, 0
   if x2 < x3 {
       dx = x3 - x2
   } else if x4 < x1 {
       dx = x1 - x4
   }
   if y2 < y3 {
       dy = y3 - y2
   } else if y4 < y1 {
       dy = y1 - y4
   }
   return dx*dx + dy*dy
}

func dist(x1, y1, x2, y2, x3, y3, x4, y4 int) float64 {
   dx, dy := 0.0, 0.0
   if x2 < x3 {
       dx = float64(x3 - x2)
   } else if x4 < x1 {
       dx = float64(x1 - x4)
   }
   if y2 < y3 {
       dy = float64(y3 - y2)
   } else if y4 < y1 {
       dy = float64(y1 - y4)
   }
   return math.Hypot(dx, dy)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var a, b, ax, ay, bx, by, n int
   fmt.Fscan(reader, &a, &b, &ax, &ay, &bx, &by, &n)
   // rectangles
   x1 := make([]int, n+2)
   y1 := make([]int, n+2)
   x2 := make([]int, n+2)
   y2 := make([]int, n+2)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x1[i], &y1[i], &x2[i], &y2[i])
       if x1[i] == x2[i] {
           if y1[i] > y2[i] {
               y1[i], y2[i] = y2[i], y1[i]
           }
       } else {
           if x1[i] > x2[i] {
               x1[i], x2[i] = x2[i], x1[i]
           }
       }
   }
   // source
   x1[n], x2[n], y1[n], y2[n] = ax, ax, ay, ay
   // target
   x1[n+1], x2[n+1], y1[n+1], y2[n+1] = bx, bx, by, by

   // BFS
   total := n + 2
   d := make([]int, total)
   for i := range d {
       d[i] = INF
   }
   queue := make([]int, 0, total)
   head := 0
   queue = append(queue, n)
   d[n] = 0
   a2 := a * a
   for head < len(queue) {
       cur := queue[head]
       head++
       if cur == n+1 {
           break
       }
       for i := 0; i < total; i++ {
           if d[i] > d[cur]+1 {
               if dist2(x1[cur], y1[cur], x2[cur], y2[cur], x1[i], y1[i], x2[i], y2[i]) <= a2 {
                   d[i] = d[cur] + 1
                   queue = append(queue, i)
               }
           }
       }
   }
   if d[n+1] == INF {
       fmt.Fprint(writer, -1)
       return
   }
   // compute answer
   hops := d[n+1]
   ans := math.Inf(1)
   for i := 0; i < n+1; i++ {
       if d[i] == hops-1 {
           cost := float64(a+b)*float64(d[i]) + dist(x1[i], y1[i], x2[i], y2[i], x1[n+1], y1[n+1], x2[n+1], y2[n+1])
           if cost < ans {
               ans = cost
           }
       }
   }
   fmt.Fprint(writer, ans)
}
