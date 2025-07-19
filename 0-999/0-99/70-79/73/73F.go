package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

// Pnt represents a point or vector in 2D
type Pnt struct { x, y float64 }

// get computes the maximum value of length * t * w / min(r, 2*pi-r)
// for t in [0,1] sampled at steps of 0.01, where r is the angle difference
// between the vector (a + dir*t) and the offset s.
func get(a, b Pnt, s, w float64) float64 {
   dir := Pnt{b.x - a.x, b.y - a.y}
   length := math.Hypot(dir.x, dir.y)
   best := 0.0
   twoPi := 2 * math.Pi
   for i := 0; i <= 100; i++ {
       t := float64(i) / 100.0
       x := a.x + dir.x*t
       y := a.y + dir.y*t
       r := math.Atan2(y, x) - s
       for r < 0 {
           r += twoPi
       }
       for r > twoPi {
           r -= twoPi
       }
       minr := math.Min(r, twoPi-r)
       val := length * t * w / minr
       if val > best {
           best = val
       }
   }
   return best
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var ax, ay, bx, by float64
   if _, err := fmt.Fscan(in, &ax, &ay); err != nil {
       return
   }
   fmt.Fscan(in, &bx, &by)
   var n int
   fmt.Fscan(in, &n)
   a := Pnt{ax, ay}
   b := Pnt{bx, by}
   v := make([]float64, 0, n+1)
   v = append(v, 0.0)
   for i := 0; i < n; i++ {
       var px, py, s, w float64
       fmt.Fscan(in, &px, &py)
       fmt.Fscan(in, &s, &w)
       p := Pnt{px, py}
       ap := Pnt{a.x - p.x, a.y - p.y}
       bp := Pnt{b.x - p.x, b.y - p.y}
       v = append(v, get(ap, bp, s, w))
   }
   sort.Float64s(v)
   for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
       v[i], v[j] = v[j], v[i]
   }
   var k int
   fmt.Fscan(in, &k)
   fmt.Fprintf(out, "%.20f\n", v[k])
}
