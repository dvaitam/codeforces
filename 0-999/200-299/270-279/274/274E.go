package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   // Obstacles on diagonals
   d1 := make(map[int][]int)
   d2 := make(map[int][]int)
   for i := 0; i < k; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       key1 := x - y
       d1[key1] = append(d1[key1], x)
       key2 := x + y
       d2[key2] = append(d2[key2], x)
   }
   for key := range d1 {
       sort.Ints(d1[key])
   }
   for key := range d2 {
       sort.Ints(d2[key])
   }
   var xs, ys int
   var dirStr string
   fmt.Fscan(in, &xs, &ys, &dirStr)
   // direction dx, dy
   dx, dy := 0, 0
   switch dirStr {
   case "NE": dx, dy = -1, 1
   case "NW": dx, dy = -1, -1
   case "SE": dx, dy = 1, 1
   case "SW": dx, dy = 1, -1
   }
   // visited states
   seen := make(map[int64]bool)
   // initial
   x, y := xs, ys
   // visited cells
   visited := make(map[int64]bool)
   // mark start
   key0 := (int64(xs) << 32) | int64(ys)
   visited[key0] = true
   ans := int64(1)
   // dir code: 0: NE,1:NW,2:SW,3:SE
   dirCode := func(dx, dy int) int {
       if dx < 0 {
           if dy > 0 {
               return 0
           }
           return 1
       }
       if dy < 0 {
           return 2
       }
       return 3
   }
   for {
       dc := dirCode(dx, dy)
       keyState := int64(x)<<32 | int64(y)<<16 | int64(dc)
       if seen[keyState] {
           break
       }
       seen[keyState] = true
       // steps to border
       var tx, ty int
       if dx > 0 {
           tx = n - x
       } else {
           tx = x - 1
       }
       if dy > 0 {
           ty = m - y
       } else {
           ty = y - 1
       }
       tBorder := tx
       if ty < tBorder {
           tBorder = ty
       }
       // steps to obstacle
       const INF = 1<<60
       tObs := INF
       if dx*dy > 0 {
           // slope +1, d1 key
           k1 := x - y
           arr := d1[k1]
           if len(arr) > 0 {
               if dx > 0 {
                   // find first arr[i] > x
                   i := sort.Search(len(arr), func(i int) bool { return arr[i] > x })
                   if i < len(arr) {
                       tObs = arr[i] - x
                   }
               } else {
                   // dx < 0, find last arr[i] < x
                   i := sort.Search(len(arr), func(i int) bool { return arr[i] >= x })
                   if i > 0 {
                       tObs = x - arr[i-1]
                   }
               }
           }
       } else {
           // slope -1, d2 key
           k2 := x + y
           arr := d2[k2]
           if len(arr) > 0 {
               if dx > 0 {
                   i := sort.Search(len(arr), func(i int) bool { return arr[i] > x })
                   if i < len(arr) {
                       tObs = arr[i] - x
                   }
               } else {
                   i := sort.Search(len(arr), func(i int) bool { return arr[i] >= x })
                   if i > 0 {
                       tObs = x - arr[i-1]
                   }
               }
           }
       }
       // move and mark each visited cell
       t := tBorder
       hitObs := false
       if int64(tObs) < int64(tBorder) {
           t = int(tObs)
           hitObs = true
       }
       // step through cells, skip marking blocked cell on obstacle hit
       for step := 1; step <= t; step++ {
           if hitObs && step == t {
               // obstacle cell, do not count
               continue
           }
           nx := x + dx*step
           ny := y + dy*step
           key := (int64(nx) << 32) | int64(ny)
           if !visited[key] {
               visited[key] = true
               ans++
           }
       }
       // update position
       x += dx * t
       y += dy * t
       // reflect
       if hitObs {
           dx = -dx
           dy = -dy
       } else {
           // border reflection
           if tx < ty {
               dx = -dx
           } else if ty < tx {
               dy = -dy
           } else {
               dx = -dx
               dy = -dy
           }
       }
   }
   fmt.Println(ans)
}
