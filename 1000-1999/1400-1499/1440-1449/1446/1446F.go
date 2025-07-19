package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

type point struct {
   x, y float64
}
type entry struct {
   ang float64
   idx int
}

var (
   n   int
   K   int64
   w   []point
   U   []entry
   BB  [][2]int
   BIT []int
   PI  = math.Acos(-1.0)
)

func normalize(t float64) float64 {
   // normalize to [-PI, PI]
   t = math.Mod(t+PI, 2*PI)
   if t < 0 {
       t += 2 * PI
   }
   return t - PI
}

// Fenwick tree on BIT with size >= c+2
func ftUpdate(i, v, size int) {
   i++
   for i < size {
       BIT[i] += v
       i += i & -i
   }
}
func ftQuery(i int) int {
   i++
   s := 0
   for i > 0 {
       s += BIT[i]
       i -= i & -i
   }
   return s
}

// Get number of pairs with distance > r
func Get(r float64) int64 {
   c := 0
   // collect intervals
   for i := 0; i < n; i++ {
       dx := w[i].x
       dy := w[i].y
       d := math.Hypot(dx, dy)
       if d <= r {
           continue
       }
       an := math.Atan2(dy, dx)
       z := math.Acos(r / d)
       b := normalize(an - z)
       e := normalize(an + z)
       U[c] = entry{ang: b, idx: i}
       c++
       U[c] = entry{ang: e, idx: i}
       c++
   }
   // sort events
   heap := U[:c]
   // simple sort by angle
   sort.Slice(heap, func(i, j int) bool {
       return heap[i].ang < heap[j].ang
   })
   // assign endpoints
   for i := 0; i < n; i++ {
       BB[i][0], BB[i][1] = -1, -1
   }
   for i := 0; i < c; i++ {
       idx := heap[i].idx
       if BB[idx][0] < 0 {
           BB[idx][0] = i
       } else {
           BB[idx][1] = i
       }
   }
   // init BIT
   size := c + 2
   for i := 0; i < size; i++ {
       BIT[i] = 0
   }
   var s int64
   total := int64(n) * int64(n-1) / 2
   // sweep
   for i := 0; i < c; i++ {
       idx := heap[i].idx
       if BB[idx][0] == i {
           s += int64(ftQuery(BB[idx][1]) - ftQuery(i))
           ftUpdate(BB[idx][1], 1, size)
       }
   }
   return total - s
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n, &K)
   w = make([]point, n)
   for i := 0; i < n; i++ {
       var xi, yi float64
       fmt.Fscan(reader, &xi, &yi)
       w[i] = point{x: xi, y: yi}
   }
   U = make([]entry, 2*n)
   BB = make([][2]int, n)
   BIT = make([]int, 2*n+5)
   low, high := 0.0, 1e5
   for it := 0; it < 60; it++ {
       mid := (low + high) * 0.5
       if Get(mid) >= K {
           high = mid
       } else {
           low = mid
       }
   }
   res := (low + high) * 0.5
   fmt.Printf("%.10f\n", res)
}
