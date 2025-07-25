package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const MOD = 998244353
var inv2 = (MOD + 1) / 2
var inv6 = 166374059 // inverse of 6 modulo MOD

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var t int
   fmt.Fscan(in, &n, &t)
   xs := make([]int, n)
   ys := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &xs[i], &ys[i])
   }
   // breakpoints
   bpMap := make(map[int]bool)
   bpMap[0] = true
   bpMap[t+1] = true
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           dx := xs[i] - xs[j]
           if dx < 0 {
               dx = -dx
           }
           dy := ys[i] - ys[j]
           if dy < 0 {
               dy = -dy
           }
           D := dx
           if dy > D {
               D = dy
           }
           // ceil(D/2)
           d0 := D / 2
           if D%2 != 0 {
               d0++
           }
           if d0 <= t {
               bpMap[d0] = true
           }
       }
   }
   bp := make([]int, 0, len(bpMap))
   for d := range bpMap {
       bp = append(bp, d)
   }
   sort.Ints(bp)
   // sample d values
   sampleMap := make(map[int]bool)
   for _, L := range bp {
       for dd := 0; dd <= 2; dd++ {
           d := L + dd
           if d >= 0 && d <= t {
               sampleMap[d] = true
           }
       }
   }
   sampleMap[t] = true
   sampleDs := make([]int, 0, len(sampleMap))
   for d := range sampleMap {
       sampleDs = append(sampleDs, d)
   }
   sort.Ints(sampleDs)
   // compute A(d) for samples
   A := make(map[int]int64)
   for _, d := range sampleDs {
       A[d] = unionArea(xs, ys, d) % MOD
   }
   // precompute prefix sums over intervals
   var sumA int64 = 0
   // iterate bp intervals
   m := len(bp)
   for k := 0; k+1 < m; k++ {
       L := bp[k]
       R := bp[k+1]
       if L > t {
           break
       }
       if R > t+1 {
           R = t+1
       }
       if L >= R {
           continue
       }
       length := R - L
       if length <= 2 {
           for d := L; d < R; d++ {
               sumA = (sumA + A[d]) % MOD
           }
       } else {
           // fit quadratic on A(L), A(L+1), A(L+2)
           y0 := A[L]
           y1 := A[L+1]
           y2 := A[L+2]
           // a = (y2 - 2*y1 + y0) / 2
           a := (y2 - 2*y1 + y0) % MOD
           if a < 0 {
               a += MOD
           }
           a = a * int64(inv2) % MOD
           // b = (y1 - y0) - a*(2L+1)
           twoL1 := int64(2*L + 1 % MOD)
           b := (y1 - y0) % MOD
           if b < 0 {
               b += MOD
           }
           b = (b - a*twoL1%MOD + MOD) % MOD
           // c = y0 - a*L^2 - b*L
           Lmod := int64(L) % MOD
           c := (y0 - a*(Lmod*Lmod%MOD)%MOD - b*Lmod%MOD) % MOD
           if c < 0 {
               c += MOD
           }
           // sum_{d=L to R-1} a d^2 + b d + c
           // sum d = S1, sum d^2 = S2
           len64 := int64(length)
           // S1 = len*L + len*(len-1)/2
           S1 := (len64*(int64(L)%MOD)%MOD + len64*(len64-1)%MOD*int64(inv2)%MOD) % MOD
           // S2 = S2(R-1) - S2(L-1)
           S2_hi := sumSq(int64(R-1))
           S2_lo := sumSq(int64(L-1))
           S2 := (S2_hi - S2_lo) % MOD
           if S2 < 0 {
               S2 += MOD
           }
           // total
           cur := (a*S2%MOD + b*S1%MOD + c*len64%MOD) % MOD
           sumA = (sumA + cur) % MOD
       }
   }
   At := A[t]
   tmod := int64(t) % MOD
   ans := (tmod*At%MOD - sumA + MOD) % MOD
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
}

// unionArea returns number of integer grid points in union of squares of Chebyshev radius d around seeds
func unionArea(xs, ys []int, d int) int64 {
   n := len(xs)
   type event struct{ x, typ, y1, y2 int }
   events := make([]event, 0, 2*n)
   for i := 0; i < n; i++ {
       x1 := xs[i] - d
       x2 := xs[i] + d
       y1 := ys[i] - d
       y2 := ys[i] + d
       events = append(events, event{x1, +1, y1, y2})
       events = append(events, event{x2 + 1, -1, y1, y2})
   }
   sort.Slice(events, func(i, j int) bool { return events[i].x < events[j].x })
   var area int64 = 0
   prevX := events[0].x
   active := make([][2]int, 0, n)
   i := 0
   for i < len(events) {
       x := events[i].x
       dx := int64(x - prevX)
       if dx > 0 && len(active) > 0 {
           // compute covered Y length
           sort.Slice(active, func(i, j int) bool { return active[i][0] < active[j][0] })
           var cov int64 = 0
           curL, curR := active[0][0], active[0][1]
           for _, iv := range active[1:] {
               if iv[0] > curR+1 {
                   cov += int64(curR - curL + 1)
                   curL, curR = iv[0], iv[1]
               } else if iv[1] > curR {
                   curR = iv[1]
               }
           }
           cov += int64(curR - curL + 1)
           area += dx * cov
       }
       // process events at x
       for i < len(events) && events[i].x == x {
           ev := events[i]
           if ev.typ > 0 {
               active = append(active, [2]int{ev.y1, ev.y2})
           } else {
               // remove one matching interval
               for j, iv := range active {
                   if iv[0] == ev.y1 && iv[1] == ev.y2 {
                       active = append(active[:j], active[j+1:]...)
                       break
                   }
               }
           }
           i++
       }
       prevX = x
   }
   return area
}

// sumSq returns sum_{i=0..n} i^2 modulo MOD
func sumSq(n int64) int64 {
   if n < 0 {
       return 0
   }
   nm := n % MOD
   nm1 := (n + 1) % MOD
   twon1 := (2*n + 1) % MOD
   return nm * nm1 % MOD * twon1 % MOD * int64(inv6) % MOD
}
