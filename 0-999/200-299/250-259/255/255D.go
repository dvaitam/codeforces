package main

import (
   "bufio"
   "fmt"
   "os"
)

func minInt64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

// S(k) = sum_{j=0..k} width(j), where width(j) = min(n, 2j+1, j+D)
// D = min(y, n-y+1)
func prefixSum(k, n, D int64) int64 {
   if k < 0 {
       return 0
   }
   // thresholds
   j1 := D - 1          // 2j+1 <= j+D for j <= j1
   j3 := n - D          // j+D <= n for j <= j3
   // sum1: j in [0..u1] width = 2j+1
   u1 := minInt64(k, j1)
   sum1 := (u1 + 1) * (u1 + 1)
   // sum2: j in [j1+1 .. u2] width = j+D
   u2 := minInt64(k, j3)
   var sum2 int64
   if u2 > j1 {
       L := j1 + 1
       R := u2
       cnt := R - L + 1
       // sum j from L to R = (R*(R+1)/2 - (L-1)*L/2)
       sumJ := (R*(R+1)/2 - (L-1)*L/2)
       sum2 = sumJ + cnt*D
   }
   // sum3: j in [j3+1 .. k] width = n
   var sum3 int64
   if k > j3 {
       cnt3 := k - j3
       sum3 = cnt3 * n
   }
   return sum1 + sum2 + sum3
}

// width at rem t: number of columns = min(n, 2t+1, t+D)
func widthAt(t, n, D int64) int64 {
   w1 := 2*t + 1
   w2 := t + D
   w := w1
   if w2 < w {
       w = w2
   }
   if n < w {
       w = n
   }
   return w
}

// f(t): total ON cells at time t
func countOn(t, n, x, y, c int64) int64 {
   D := minInt64(y, n-y+1)
   // full width at rem t (row x)
   w0 := widthAt(t, n, D)
   // rows above and below
   up := t
   if up > x-1 {
       up = x - 1
   }
   down := t
   if down > n-x {
       down = n - x
   }
   // S(t-1)
   S_t1 := prefixSum(t-1, n, D)
   S_up := prefixSum(t-up-1, n, D)
   S_down := prefixSum(t-down-1, n, D)
   // sum of widths for rows above and below
   sumRows := (S_t1 - S_up) + (S_t1 - S_down)
   return w0 + sumRows
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, x, y, c int64
   if _, err := fmt.Fscan(in, &n, &x, &y, &c); err != nil {
       return
   }
   // maximum needed t is max distance to any corner
   var maxDist int64
   corners := []struct{ dx, dy int64 }{{x - 1, y - 1}, {x - 1, n - y}, {n - x, y - 1}, {n - x, n - y}}
   for _, d := range corners {
       s := d.dx + d.dy
       if s > maxDist {
           maxDist = s
       }
   }
   lo, hi := int64(-1), maxDist
   for lo+1 < hi {
       mid := (lo + hi) / 2
       if countOn(mid, n, x, y, c) >= c {
           hi = mid
       } else {
           lo = mid
       }
   }
   // check hi
   if hi < 0 {
       hi = 0
   }
   fmt.Println(hi)
}
