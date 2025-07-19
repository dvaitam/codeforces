package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// PointIdx holds a point's coordinates and original index
type PointIdx struct {
   x, y int64
   idx  int
}

// cross returns the cross product of vectors (x1,y1) and (x2,y2)
func cross(x1, y1, x2, y2 int64) int64 {
   return x1*y2 - y1*x2
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   // use 1-based indexing for convenience, length 5 covers up to P[4]
   P := make([]PointIdx, 5)
   for i := 1; i <= 3; i++ {
       fmt.Fscan(reader, &P[i].x, &P[i].y)
       P[i].idx = i
   }
   // sort initial three points by x, then y
   sort.Slice(P[1:4], func(i, j int) bool {
       a := P[i+1]
       b := P[j+1]
       if a.x != b.x {
           return a.x < b.x
       }
       return a.y < b.y
   })

   // process remaining points
   for i := 4; i <= n; i++ {
       var x, y int64
       fmt.Fscan(reader, &x, &y)
       ax, ay := P[1].x, P[1].y
       bx, by := P[2].x, P[2].y
       cx, cy := P[3].x, P[3].y

       // check if initial three are collinear
       if cross(bx-ax, by-ay, cx-ax, cy-ay) == 0 {
           // if new point is also collinear with line P1-P2
           if cross(bx-ax, by-ay, x-ax, y-ay) == 0 {
               P[4] = PointIdx{x: x, y: y, idx: i}
               // sort four collinear points
               sort.Slice(P[1:5], func(ii, jj int) bool {
                   aa := P[ii+1]
                   bb := P[jj+1]
                   if aa.x != bb.x {
                       return aa.x < bb.x
                   }
                   return aa.y < bb.y
               })
           } else {
               // replace third point
               P[3] = PointIdx{x: x, y: y, idx: i}
           }
           continue
       }
       // compute cross products for triangle orientation
       f1 := cross(bx-ax, by-ay, x-ax, y-ay)
       f2 := cross(cx-bx, cy-by, x-bx, y-by)
       f3 := cross(ax-cx, ay-cy, x-cx, y-cy)
       // if new point lies inside or on edges of triangle
       if (f1 <= 0 && f2 <= 0 && f3 <= 0) || (f1 >= 0 && f2 >= 0 && f3 >= 0) {
           if f1 == 0 {
               P[2] = PointIdx{x: x, y: y, idx: i}
           } else {
               P[3] = PointIdx{x: x, y: y, idx: i}
           }
       }
   }
   // output the indices of the three non-collinear points
   fmt.Fprintf(writer, "%d %d %d", P[1].idx, P[2].idx, P[3].idx)
}
