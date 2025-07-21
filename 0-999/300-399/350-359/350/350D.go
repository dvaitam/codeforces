package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   if a < 0 {
       return -a
   }
   return a
}

func gcd3(a, b, c int64) int64 {
   return gcd(gcd(a, b), c)
}

type lineKey struct {
   a, b, c int64
}

type lineGroup struct {
   starts []int64
   ends   []int64
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   // map lineKey to segments covering intervals
   groups := make(map[lineKey]*lineGroup, n)
   for i := 0; i < n; i++ {
       var x1, y1, x2, y2 int64
       fmt.Fscan(reader, &x1, &y1, &x2, &y2)
       dx := x2 - x1
       dy := y2 - y1
       g := gcd(dx, dy)
       dx /= g
       dy /= g
       if dx < 0 || (dx == 0 && dy < 0) {
           dx = -dx
           dy = -dy
       }
       a := -dy
       b := dx
       c2 := 2*(a*x1 + b*y1)
       g2 := gcd3(a, b, c2)
       a /= g2
       b /= g2
       c2 /= g2
       if a < 0 || (a == 0 && b < 0) {
           a = -a
           b = -b
           c2 = -c2
       }
       key := lineKey{a, b, c2}
       grp, ok := groups[key]
       if !ok {
           grp = &lineGroup{}
           groups[key] = grp
       }
       t1 := 2*(dx*x1 + dy*y1)
       t2 := 2*(dx*x2 + dy*y2)
       if t1 <= t2 {
           grp.starts = append(grp.starts, t1)
           grp.ends = append(grp.ends, t2)
       } else {
           grp.starts = append(grp.starts, t2)
           grp.ends = append(grp.ends, t1)
       }
   }
   // sort intervals
   for _, grp := range groups {
       sort.Slice(grp.starts, func(i, j int) bool { return grp.starts[i] < grp.starts[j] })
       sort.Slice(grp.ends, func(i, j int) bool { return grp.ends[i] < grp.ends[j] })
   }
   // read circles
   type circle struct{ x, y, r int64 }
   byRadius := make(map[int64][]circle)
   circles := make([]circle, m)
   for i := 0; i < m; i++ {
       var x, y, r int64
       fmt.Fscan(reader, &x, &y, &r)
       circles[i] = circle{x, y, r}
       byRadius[r] = append(byRadius[r], circles[i])
   }
   var ans int64
   // process pairs
   for r, cs := range byRadius {
       k := len(cs)
       if k < 2 {
           continue
       }
       // for each pair
       lim := 4 * r * r
       for i := 0; i < k; i++ {
           xi, yi := cs[i].x, cs[i].y
           for j := i + 1; j < k; j++ {
               xj, yj := cs[j].x, cs[j].y
               dx := xj - xi
               dy := yj - yi
               if dx*dx+dy*dy <= lim {
                   continue
               }
               // perpendicular bisector direction
               dxd := -dy
               dyd := dx
               g := gcd(dxd, dyd)
               dxd /= g
               dyd /= g
               if dxd < 0 || (dxd == 0 && dyd < 0) {
                   dxd = -dxd
                   dyd = -dyd
               }
               // line coefficients
               a := -dyd
               b := dxd
               c2 := a*(xi+xj) + b*(yi+yj)
               g2 := gcd3(a, b, c2)
               a /= g2
               b /= g2
               c2 /= g2
               if a < 0 || (a == 0 && b < 0) {
                   a = -a
                   b = -b
                   c2 = -c2
               }
               key := lineKey{a, b, c2}
               grp, ok := groups[key]
               if !ok {
                   continue
               }
               // midpoint projection
               tmid := dxd*(xi+xj) + dyd*(yi+yj)
               // count segments covering tmid
               sCnt := sort.Search(len(grp.starts), func(i int) bool { return grp.starts[i] > tmid })
               eCnt := sort.Search(len(grp.ends), func(i int) bool { return grp.ends[i] >= tmid })
               ans += int64(sCnt - eCnt)
           }
       }
   }
   fmt.Fprint(writer, ans)
}
