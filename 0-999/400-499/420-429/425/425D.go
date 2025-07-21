package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   pts := make([][2]int, 0, n)
   xsAtY := make(map[int][]int, n)
   ysAtX := make(map[int][]int, n)
   pointSet := make(map[int64]struct{}, n*2)

   for i := 0; i < n; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       pts = append(pts, [2]int{x, y})
       xsAtY[y] = append(xsAtY[y], x)
       ysAtX[x] = append(ysAtX[x], y)
       key := (int64(x) << 32) | int64(y)
       pointSet[key] = struct{}{}
   }
   // sort coordinate lists
   for y, xs := range xsAtY {
       sort.Ints(xs)
       xsAtY[y] = xs
   }
   for x, ys := range ysAtX {
       sort.Ints(ys)
       ysAtX[x] = ys
   }

   var result int64 = 0
   // for each point as bottom-left
   for _, p := range pts {
       x, y := p[0], p[1]
       xs := xsAtY[y]
       ys := ysAtX[x]
       // choose smaller list to iterate
       if len(xs) < len(ys) {
           for _, x2 := range xs {
               if x2 <= x {
                   continue
               }
               d := x2 - x
               // check upper two points
               key1 := (int64(x) << 32) | int64(y+d)
               key2 := (int64(x2) << 32) | int64(y+d)
               if _, ok1 := pointSet[key1]; ok1 {
                   if _, ok2 := pointSet[key2]; ok2 {
                       result++
                   }
               }
           }
       } else {
           for _, y2 := range ys {
               if y2 <= y {
                   continue
               }
               d := y2 - y
               // check right two points
               key1 := (int64(x+d) << 32) | int64(y)
               key2 := (int64(x+d) << 32) | int64(y2)
               if _, ok1 := pointSet[key1]; ok1 {
                   if _, ok2 := pointSet[key2]; ok2 {
                       result++
                   }
               }
           }
       }
   }
   fmt.Fprint(writer, result)
}
