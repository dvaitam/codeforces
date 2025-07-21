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

func absInt(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   points := make([]Point, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &points[i].x, &points[i].y)
   }
   // indices sorted by x and by y
   sx := make([]int, n)
   sy := make([]int, n)
   for i := 0; i < n; i++ {
       sx[i] = i
       sy[i] = i
   }
   sort.Slice(sx, func(i, j int) bool {
       if points[sx[i]].x != points[sx[j]].x {
           return points[sx[i]].x < points[sx[j]].x
       }
       return points[sx[i]].y < points[sx[j]].y
   })
   sort.Slice(sy, func(i, j int) bool {
       if points[sy[i]].y != points[sy[j]].y {
           return points[sy[i]].y < points[sy[j]].y
       }
       return points[sy[i]].x < points[sy[j]].x
   })
   type task struct{ sx, sy []int }
   stack := make([]task, 0, 16)
   stack = append(stack, task{sx: sx, sy: sy})
   ans := 0
   for len(stack) > 0 {
       // pop
       t := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       m := len(t.sx)
       if m == 0 {
           continue
       }
       // try split on x
       // build unique xs and counts
       ux := make([]int64, 0, m)
       cx := make([]int, 0, m)
       // build ux, cx by scanning t.sx
       for i, idx := range t.sx {
           if i == 0 || points[idx].x != ux[len(ux)-1] {
               ux = append(ux, points[idx].x)
               cx = append(cx, 1)
           } else {
               cx[len(cx)-1]++
           }
       }
       // prefix sums of counts
       px := make([]int, len(cx))
       for i := range cx {
           if i == 0 {
               px[i] = cx[i]
           } else {
               px[i] = px[i-1] + cx[i]
           }
       }
       // find best gap in x
       bestDiff := m + 1
       bestK := -1
       for k := 0; k+1 < len(ux); k++ {
           if ux[k+1] > ux[k]+1 {
               left := px[k]
               diff := absInt(left - (m - left))
               if diff < bestDiff {
                   bestDiff = diff
                   bestK = k
               }
           }
       }
       if bestK >= 0 {
           // split on x at ux[bestK]
           // left in t.sx are first px[bestK]
           lx := t.sx[:px[bestK]]
           rx := t.sx[px[bestK]:]
           // partition sy by x
           ly := make([]int, 0, len(lx))
           ry := make([]int, 0, len(rx))
           xThr := ux[bestK]
           for _, idx := range t.sy {
               if points[idx].x <= xThr {
                   ly = append(ly, idx)
               } else {
                   ry = append(ry, idx)
               }
           }
           stack = append(stack, task{sx: rx, sy: ry})
           stack = append(stack, task{sx: lx, sy: ly})
           continue
       }
       // try split on y
       uy := make([]int64, 0, m)
       cy := make([]int, 0, m)
       for i, idx := range t.sy {
           if i == 0 || points[idx].y != uy[len(uy)-1] {
               uy = append(uy, points[idx].y)
               cy = append(cy, 1)
           } else {
               cy[len(cy)-1]++
           }
       }
       py := make([]int, len(cy))
       for i := range cy {
           if i == 0 {
               py[i] = cy[i]
           } else {
               py[i] = py[i-1] + cy[i]
           }
       }
       bestDiff = m + 1
       bestK = -1
       for k := 0; k+1 < len(uy); k++ {
           if uy[k+1] > uy[k]+1 {
               left := py[k]
               diff := absInt(left - (m - left))
               if diff < bestDiff {
                   bestDiff = diff
                   bestK = k
               }
           }
       }
       if bestK >= 0 {
           // split on y at uy[bestK]
           ly := t.sy[:py[bestK]]
           ry := t.sy[py[bestK]:]
           // partition sx by y
           lx := make([]int, 0, len(ly))
           rx := make([]int, 0, len(ry))
           yThr := uy[bestK]
           for _, idx := range t.sx {
               if points[idx].y <= yThr {
                   lx = append(lx, idx)
               } else {
                   rx = append(rx, idx)
               }
           }
           stack = append(stack, task{sx: rx, sy: ry})
           stack = append(stack, task{sx: lx, sy: ly})
           continue
       }
       // no splits possible
       ans++
   }
   // output result
   fmt.Println(ans)
