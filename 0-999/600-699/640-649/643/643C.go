package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

// Line represents y = m*x + b
type Line struct {
   m, b float64
}

// CHT for minimum query, slopes inserted in decreasing order, queries with increasing x
type CHT struct {
   lines []Line
   head  int
}

// isBad checks whether line l2 is unnecessary between l1 and l3
func isBad(l1, l2, l3 Line) bool {
   // intersection x-coordinate of l1 and l2: (b2-b1)/(m1-m2)
   // l2 and l3: (b3-b2)/(m2-m3)
   return (l2.b-l1.b)*(l2.m-l3.m) >= (l3.b-l2.b)*(l1.m-l2.m)
}

// AddLine adds a new line with slope m and intercept b
func (h *CHT) AddLine(m, b float64) {
   h.lines = append(h.lines, Line{m: m, b: b})
   // remove middle lines while they are bad
   for len(h.lines)-h.head >= 3 {
       n := len(h.lines)
       if isBad(h.lines[n-3], h.lines[n-2], h.lines[n-1]) {
           // remove second last
           h.lines = append(h.lines[:n-2], h.lines[n-1])
       } else {
           break
       }
   }
}

// Query returns min value at x
func (h *CHT) Query(x float64) float64 {
   // move head while next line gives smaller value
   for h.head+1 < len(h.lines) && h.lines[h.head].m*x+h.lines[h.head].b >= h.lines[h.head+1].m*x+h.lines[h.head+1].b {
       h.head++
   }
   return h.lines[h.head].m*x + h.lines[h.head].b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   t := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &t[i])
   }
   // prefix sums
   S := make([]float64, n+1)
   P1 := make([]float64, n+1)
   P2 := make([]float64, n+1)
   for i := 1; i <= n; i++ {
       S[i] = S[i-1] + float64(t[i])
       P1[i] = P1[i-1] + S[i]/float64(t[i])
       P2[i] = P2[i-1] + 1.0/float64(t[i])
   }
   inf := math.Inf(1)
   dpPrev := make([]float64, n+1)
   dpCurr := make([]float64, n+1)
   dpPrev[0] = 0.0
   for i := 1; i <= n; i++ {
       dpPrev[i] = inf
   }
   // DP for c = 1..k
   for c := 1; c <= k; c++ {
       // init CHT
       var cht CHT
       // init dpCurr
       dpCurr[0] = inf
       for i := 1; i <= n; i++ {
           dpCurr[i] = inf
       }
       for i := 1; i <= n; i++ {
           j := i - 1
           // compute line from j
           Cj := dpPrev[j] - P1[j] + S[j]*P2[j]
           // slope = -S[j]
           cht.AddLine(-S[j], Cj)
           // query for i
           val := cht.Query(P2[i])
           dpCurr[i] = P1[i] + val
       }
       // swap
       dpPrev, dpCurr = dpCurr, dpPrev
   }
   // answer dpPrev[n]
   ans := dpPrev[n]
   // print with precision
   fmt.Printf("%.10f\n", ans)
}
