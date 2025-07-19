package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

// Point holds absolute coordinates and original index
type Point struct {
   X, Y float64
   ind  int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   xs := make([]int, n)
   ys := make([]int, n)
   absMap := make(map[[2]int]int)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &xs[i], &ys[i])
       ax := xs[i]
       if ax < 0 {
           ax = -ax
       }
       ay := ys[i]
       if ay < 0 {
           ay = -ay
       }
       key := [2]int{ax, ay}
       if j, ok := absMap[key]; ok {
           // duplicate absolute pair, print and exit
           cat := category(xs[j], ys[j], xs[i], ys[i])
           fmt.Fprintf(writer, "%d %d %d 1\n", j+1, cat, i+1)
           return
       }
       absMap[key] = i
   }
   // prepare points with absolute coords
   v := make([]Point, n)
   for i := 0; i < n; i++ {
       ax := xs[i]
       if ax < 0 {
           ax = -ax
       }
       ay := ys[i]
       if ay < 0 {
           ay = -ay
       }
       v[i] = Point{float64(ax), float64(ay), i}
   }
   // sort by X then Y
   sort.Slice(v, func(i, j int) bool {
       if v[i].X != v[j].X {
           return v[i].X < v[j].X
       }
       return v[i].Y < v[j].Y
   })
   liveY := make([]Point, 0, 16)
   window := make([]Point, 0, 16)
   mn := math.Inf(1)
   var aIdx, bIdx int
   for _, p := range v {
       // remove points too far in X
       for len(window) > 0 && p.X-window[0].X > mn {
           old := window[0]
           window = window[1:]
           // remove old from liveY
           idx := sort.Search(len(liveY), func(i int) bool {
               if liveY[i].Y != old.Y {
                   return liveY[i].Y >= old.Y
               }
               return liveY[i].X >= old.X
           })
           for idx < len(liveY) && liveY[idx].Y == old.Y && liveY[idx].X == old.X {
               if liveY[idx].ind == old.ind {
                   liveY = append(liveY[:idx], liveY[idx+1:]...)
                   break
               }
               idx++
           }
       }
       // search candidates within Y range
       lowY := p.Y - mn
       highY := p.Y + mn
       l := sort.Search(len(liveY), func(i int) bool { return liveY[i].Y >= lowY })
       for j := l; j < len(liveY); j++ {
           q := liveY[j]
           if q.Y > highY {
               break
           }
           dx := p.X - q.X
           dy := p.Y - q.Y
           d := math.Hypot(dx, dy)
           if d < mn {
               mn = d
               aIdx = p.ind
               bIdx = q.ind
           }
       }
       // add p to window and liveY
       window = append(window, p)
       idx := sort.Search(len(liveY), func(i int) bool { return liveY[i].Y >= p.Y })
       liveY = append(liveY, Point{})
       copy(liveY[idx+1:], liveY[idx:])
       liveY[idx] = p
   }
   // output result
   cat := category(xs[aIdx], ys[aIdx], xs[bIdx], ys[bIdx])
   fmt.Fprintf(writer, "%d %d %d 1\n", aIdx+1, cat, bIdx+1)
}

// category determines the relation code based on sign of coordinates
func category(x1, y1, x2, y2 int) int {
   xx := (x1 < 0 && x2 < 0) || (x1 > 0 && x2 > 0)
   yy := (y1 < 0 && y2 < 0) || (y1 > 0 && y2 > 0)
   if !xx && !yy {
       return 1
   } else if xx && !yy {
       return 2
   } else if !xx && yy {
       return 3
   }
   return 4
}
