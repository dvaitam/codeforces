package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const INF = 1e20

type Particle struct {
   x, v float64
}

func solve(p []Particle, a, b, c int) float64 {
   if a == b || b == c {
       return INF
   }
   lo, hi := 0.0, 1e9
   for rep := 0; rep < 60; rep++ {
       mid := (lo + hi) / 2.0
       f := -INF
       for i := a; i < b; i++ {
           if v := p[i].x + p[i].v*mid; v > f {
               f = v
           }
       }
       ok := true
       for i := b; i < c; i++ {
           if p[i].x + p[i].v*mid <= f {
               ok = false
               break
           }
       }
       if ok {
           lo = mid
       } else {
           hi = mid
       }
   }
   return hi
}

func main() {
   rdr := bufio.NewReader(os.Stdin)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   var n int
   if _, err := fmt.Fscan(rdr, &n); err != nil {
       return
   }
   p := make([]Particle, n)
   for i := 0; i < n; i++ {
       var xi, vi int
       fmt.Fscan(rdr, &xi, &vi)
       p[i].x = float64(xi)
       p[i].v = float64(vi)
   }
   sort.Slice(p, func(i, j int) bool { return p[i].x < p[j].x })
   ans := INF
   for i := 0; i < n; {
       j := i
       for j < n && p[j].v > 0 {
           j++
       }
       k := j
       for k < n && p[k].v < 0 {
           k++
       }
       t := solve(p, i, j, k)
       if t < ans {
           ans = t
       }
       i = k
   }
   if ans >= INF {
       fmt.Fprintln(w, -1)
   } else {
       fmt.Fprintf(w, "%.12f\n", ans)
   }
}
