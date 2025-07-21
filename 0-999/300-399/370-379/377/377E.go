package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Line represents y = m*x + b
type Line struct {
   m, b float64
}

// ConvexHull for minimum queries. Slopes (m) must be added in decreasing order.
type ConvexHull struct {
   lines []Line     // lines[i]
   xs    []float64  // xs[i]: start x-coordinate where lines[i] is optimal
}

// Empty returns true if hull has no lines
func (ch *ConvexHull) Empty() bool { return len(ch.lines) == 0 }

// addLine adds a new line y = m*x + b. m should be <= previous m.
func (ch *ConvexHull) addLine(m, b float64) {
   // remove last while intersection is not increasing
   for len(ch.lines) > 0 {
       last := ch.lines[len(ch.lines)-1]
       // parallel lines: keep one with smaller b
       if m == last.m {
           if b >= last.b {
               return
           }
           ch.lines = ch.lines[:len(ch.lines)-1]
           ch.xs = ch.xs[:len(ch.xs)-1]
           continue
       }
       // intersection x-coordinate with last line
       x := (last.b - b) / (m - last.m)
       if len(ch.xs) > 0 && x <= ch.xs[len(ch.xs)-1] {
           ch.lines = ch.lines[:len(ch.lines)-1]
           ch.xs = ch.xs[:len(ch.xs)-1]
           continue
       }
       ch.xs = append(ch.xs, x)
       break
   }
   if len(ch.lines) == 0 {
       // valid from -infinity
       ch.xs = append(ch.xs, -1e300)
   }
   ch.lines = append(ch.lines, Line{m: m, b: b})
}

// query returns minimum y-value at given x
func (ch *ConvexHull) query(x float64) float64 {
   // binary search highest i such that xs[i] <= x
   lo, hi := 0, len(ch.lines)-1
   for lo < hi {
       mid := (lo + hi + 1) / 2
       if ch.xs[mid] <= x {
           lo = mid
       } else {
           hi = mid - 1
       }
   }
   ln := ch.lines[lo]
   return ln.m*x + ln.b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var s int64
   fmt.Fscan(in, &n, &s)
   type Building struct{ v, c int }
   bs := make([]Building, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &bs[i].v, &bs[i].c)
   }
   // For equal v, keep minimal c
   best := make(map[int]int)
   for _, b := range bs {
       if pc, ok := best[b.v]; !ok || b.c < pc {
           best[b.v] = b.c
       }
   }
   vs := make([]int, 0, len(best))
   for v := range best {
       vs = append(vs, v)
   }
   sort.Ints(vs)
   filtered := make([]Building, len(vs))
   for i, v := range vs {
       filtered[i] = Building{v: v, c: best[v]}
   }
   // dp[i]: minimal time to acquire building i (in continuous)
   const INF = 1e300
   dp := make([]float64, len(filtered))
   ch := ConvexHull{}
   // process by increasing v
   for i, b := range filtered {
       if b.c == 0 {
           dp[i] = 0
       } else if ch.Empty() {
           dp[i] = INF
       } else {
           dp[i] = ch.query(float64(b.c))
       }
       if dp[i] < INF {
           // add line for this building: time + c/v -> slope=1/v, intercept=dp
           ch.addLine(1.0/float64(b.v), dp[i])
       }
   }
   // compute answer
   res := 1e300
   for i, b := range filtered {
       if dp[i] >= INF {
           continue
       }
       t := dp[i] + float64(s)/float64(b.v)
       if t < res {
           res = t
       }
   }
   // ceil result
   ans := int64(res)
   if float64(ans) < res-1e-12 {
       ans++
   }
   fmt.Println(ans)
}
