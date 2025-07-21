package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   xs := make([]int, n+1)
   for i := 0; i <= n; i++ {
       fmt.Fscan(in, &xs[i])
   }
   var py int
   fmt.Fscan(in, &py)
   // axis cities x coordinates
   axis := make([]int, n)
   copy(axis, xs[:n])
   sort.Ints(axis)
   L := float64(axis[0])
   R := float64(axis[n-1])
   Dline := R - L
   // off-axis city
   px := float64(xs[n])
   pyf := float64(py)
   // helper for distance
   dist := func(x1, y1, x2, y2 float64) float64 {
       return math.Hypot(x1-x2, y1-y2)
   }
   // for small n<=2, brute-force permutations
   if n <= 2 {
       // build nodes
       m := n + 1
       xsF := make([]float64, m)
       ysF := make([]float64, m)
       for i := 0; i < n; i++ {
           xsF[i] = float64(xs[i])
           ysF[i] = 0
       }
       xsF[n] = float64(xs[n]); ysF[n] = float64(py)
       start := k - 1
       // list of other indices
       others := make([]int, 0, m-1)
       for i := 0; i < m; i++ {
           if i != start {
               others = append(others, i)
           }
       }
       best := math.Inf(1)
       // permute others (length <=2)
       var permute func([]int, int)
       permute = func(a []int, l int) {
           if l == len(a)-1 {
               // compute path
               cur := float64(0)
               prev := start
               for _, idx := range a {
                   cur += dist(xsF[prev], ysF[prev], xsF[idx], ysF[idx])
                   prev = idx
               }
               if cur < best {
                   best = cur
               }
               return
           }
           for i := l; i < len(a); i++ {
               a[l], a[i] = a[i], a[l]
               permute(a, l+1)
               a[l], a[i] = a[i], a[l]
           }
       }
       permute(others, 0)
       fmt.Printf("%.10f\n", best)
       return
   }
   ans := math.Inf(1)
   // find closest axis city to off-axis for middle detour
   // find idx in axis where axis[idx] >= px
   idx := sort.Search(n, func(i int) bool { return float64(axis[i]) >= px })
   // candidate positions
   var dv float64 = math.Inf(1)
   for _, j := range []int{idx - 1, idx} {
       if j >= 0 && j < n {
           d := dist(px, pyf, float64(axis[j]), 0)
           if d < dv {
               dv = d
           }
       }
   }
   if k == n+1 {
       // start at off-axis
       // go to nearest axis end, traverse line
       d1 := dist(px, pyf, L, 0)
       d2 := dist(px, pyf, R, 0)
       ans = Dline + math.Min(d1, d2)
   } else {
       // start at axis city
       Sx := float64(xs[k-1])
       // two directions: start covering from L to R, or R to L
       for dir := 0; dir < 2; dir++ {
           var E1, E2 float64
           if dir == 0 {
               E1, E2 = L, R
           } else {
               E1, E2 = R, L
           }
           // cost to reach first end
           cS := math.Abs(Sx - E1)
           // baseline cover cost
           base := cS + Dline
           // 1. detour in middle (go to off-axis and return)
           ans = math.Min(ans, base + 2*dv)
           // 2. go to off-axis first (start->P->E1->E2)
           ans = math.Min(ans, dist(Sx, 0, px, pyf) + dist(px, pyf, E1, 0) + Dline)
           // 3. go to off-axis last (cover axis then P)
           ans = math.Min(ans, base + dist(px, pyf, E2, 0))
       }
   }
   // output with precision
   fmt.Printf("%.10f\n", ans)
}
