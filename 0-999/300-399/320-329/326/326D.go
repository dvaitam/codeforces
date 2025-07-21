package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// DSU structure
type DSU struct {
   p, r []int
}

func newDSU(n int) *DSU {
   p := make([]int, n)
   r := make([]int, n)
   for i := 0; i < n; i++ {
       p[i] = i
   }
   return &DSU{p: p, r: r}
}

func (d *DSU) find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.find(d.p[x])
   }
   return d.p[x]
}

func (d *DSU) union(x, y int) {
   rx := d.find(x)
   ry := d.find(y)
   if rx == ry {
       return
   }
   if d.r[rx] < d.r[ry] {
       d.p[rx] = ry
   } else if d.r[ry] < d.r[rx] {
       d.p[ry] = rx
   } else {
       d.p[ry] = rx
       d.r[rx]++
   }
}

// Edge represents an interval [low, high) on an axis for a rectangle side
type Edge struct {
   low, high, id int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   rects := make([][4]int, n)
   left := make(map[int][]Edge)
   right := make(map[int][]Edge)
   bottom := make(map[int][]Edge)
   top := make(map[int][]Edge)
   for i := 0; i < n; i++ {
       x1, y1, x2, y2 := 0, 0, 0, 0
       fmt.Fscan(in, &x1, &y1, &x2, &y2)
       rects[i] = [4]int{x1, y1, x2, y2}
       left[x1] = append(left[x1], Edge{y1, y2, i})
       right[x2] = append(right[x2], Edge{y1, y2, i})
       bottom[y1] = append(bottom[y1], Edge{x1, x2, i})
       top[y2] = append(top[y2], Edge{x1, x2, i})
   }
   dsu := newDSU(n)
   // helper to process adjacency at each coordinate
   // process adjacency: for each coord in a and b, union overlapping edges
   process := func(a, b map[int][]Edge) {
       for coord, edgesA := range a {
           edgesB, ok := b[coord]
           if !ok {
               continue
           }
           // sort by low interval endpoint
           sort.Slice(edgesA, func(i, j int) bool { return edgesA[i].low < edgesA[j].low })
           sort.Slice(edgesB, func(i, j int) bool { return edgesB[i].low < edgesB[j].low })
           j := 0
           for _, ea := range edgesA {
               for j < len(edgesB) && edgesB[j].high <= ea.low {
                   j++
               }
               for k := j; k < len(edgesB) && edgesB[k].low < ea.high; k++ {
                   dsu.union(ea.id, edgesB[k].id)
               }
           }
       }
   }
   // import sort functionality
   // process vertical adjacency (left-right)
   process(left, right)
   process(right, left)
   // process horizontal adjacency (bottom-top)
   process(bottom, top)
   process(top, bottom)

   // component stats
   minX := make([]int, n)
   maxX := make([]int, n)
   minY := make([]int, n)
   maxY := make([]int, n)
   areaSum := make([]int, n)
   for i := 0; i < n; i++ {
       minX[i] = 1<<30
       minY[i] = 1<<30
   }
   for i := 0; i < n; i++ {
       r := dsu.find(i)
       x1, y1, x2, y2 := rects[i][0], rects[i][1], rects[i][2], rects[i][3]
       if x1 < minX[r] {
           minX[r] = x1
       }
       if y1 < minY[r] {
           minY[r] = y1
       }
       if x2 > maxX[r] {
           maxX[r] = x2
       }
       if y2 > maxY[r] {
           maxY[r] = y2
       }
       areaSum[r] += (x2 - x1) * (y2 - y1)
   }
   // find valid component
   for i := 0; i < n; i++ {
       if dsu.find(i) != i {
           continue
       }
       dx := maxX[i] - minX[i]
       dy := maxY[i] - minY[i]
       if dx > 0 && dx == dy && areaSum[i] == dx*dy {
           // collect ids
           res := make([]int, 0)
           for j := 0; j < n; j++ {
               if dsu.find(j) == i {
                   res = append(res, j+1)
               }
           }
           fmt.Fprintln(out, "YES")
           fmt.Fprintln(out, len(res))
           for k, id := range res {
               if k > 0 {
                   out.WriteByte(' ')
               }
               fmt.Fprint(out, id)
           }
           out.WriteByte('\n')
           return
       }
   }
   fmt.Fprintln(out, "NO")
}
