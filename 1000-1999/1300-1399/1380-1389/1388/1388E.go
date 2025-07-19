package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

type line struct{ k, m float64 }
// cht for maximum query, slopes added in increasing order, queries in increasing x
type cht struct{
   lines []line
   ptr   int
}
// isBad checks whether middle line b is unnecessary between a and c
func isBad(a, b, c line) bool {
   // intersection x of a&b: (b.m-a.m)/(a.k-b.k)
   // intersection x of b&c: (c.m-b.m)/(b.k-c.k)
   // b is bad if x_ab >= x_bc
   return (b.m-a.m)*(b.k-c.k) >= (c.m-b.m)*(a.k-b.k)
}
// add line with slope k, intercept m
func (h *cht) add(k, m float64) {
   ln := line{k: k, m: m}
   // remove last while bad
   l := h.lines
   for len(l) >= 2 && isBad(l[len(l)-2], l[len(l)-1], ln) {
       l = l[:len(l)-1]
   }
   h.lines = append(l, ln)
}
// query maximum value at x, queries must be non-decreasing
func (h *cht) query(x float64) float64 {
   // move ptr while next is better
   for h.ptr+1 < len(h.lines) &&
       h.lines[h.ptr].k*x+h.lines[h.ptr].m <= h.lines[h.ptr+1].k*x+h.lines[h.ptr+1].m {
       h.ptr++
   }
   return h.lines[h.ptr].k*x + h.lines[h.ptr].m
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   xl := make([]float64, n)
   xr := make([]float64, n)
   y := make([]float64, n)
   for i := 0; i < n; i++ {
       var xi, xj, yi float64
       fmt.Fscan(reader, &xi, &xj, &yi)
       xl[i], xr[i], y[i] = xi, xj, yi
   }
   // build segments of x where width function changes
   segs := make([][2]float64, 0)
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           if y[i] == y[j] {
               continue
           }
           a := (xr[j] - xl[i]) / (y[i] - y[j])
           b := (xl[j] - xr[i]) / (y[i] - y[j])
           if a > b {
               a, b = b, a
           }
           segs = append(segs, [2]float64{a, b})
       }
   }
   sort.Slice(segs, func(i, j int) bool { return segs[i][0] < segs[j][0] })
   mr := make([][2]float64, 0)
   for _, s := range segs {
       if len(mr) == 0 || mr[len(mr)-1][1] <= s[0] {
           mr = append(mr, s)
       } else if mr[len(mr)-1][1] < s[1] {
           mr[len(mr)-1][1] = s[1]
       }
   }
   if len(mr) == 0 {
       mr = append(mr, [2]float64{0, 0})
   }
   // prepare convex hulls
   minMap := make(map[float64]float64)
   maxMap := make(map[float64]float64)
   for i := 0; i < n; i++ {
       k1 := -y[i]
       m1 := -xl[i]
       if v, ok := minMap[k1]; !ok || v < m1 {
           minMap[k1] = m1
       }
       m2 := -xr[i]
       if v, ok := minMap[k1]; !ok || v < m2 {
           minMap[k1] = m2
       }
       k2 := y[i]
       m3 := xl[i]
       if v, ok := maxMap[k2]; !ok || v < m3 {
           maxMap[k2] = m3
       }
       m4 := xr[i]
       if v, ok := maxMap[k2]; !ok || v < m4 {
           maxMap[k2] = m4
       }
   }
   minLines := make([]line, 0, len(minMap))
   for k, m := range minMap {
       minLines = append(minLines, line{k: k, m: m})
   }
   maxLines := make([]line, 0, len(maxMap))
   for k, m := range maxMap {
       maxLines = append(maxLines, line{k: k, m: m})
   }
   sort.Slice(minLines, func(i, j int) bool { return minLines[i].k < minLines[j].k })
   sort.Slice(maxLines, func(i, j int) bool { return maxLines[i].k < maxLines[j].k })
   hmin := &cht{lines: make([]line, 0, len(minLines))}
   for _, ln := range minLines {
       hmin.add(ln.k, ln.m)
   }
   hmax := &cht{lines: make([]line, 0, len(maxLines))}
   for _, ln := range maxLines {
       hmax.add(ln.k, ln.m)
   }
   ans := math.Inf(1)
   for _, seg := range mr {
       for _, x := range []float64{seg[0], seg[1]} {
           v := -hmin.query(x) + hmax.query(x)
           if v < ans {
               ans = v
           }
       }
   }
   fmt.Fprintf(writer, "%.15f", ans)
}
