package main

import (
   "bufio"
   "fmt"
   "os"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

func f(x int64) int {
   res := 0
   for x > 0 {
       d := x % 10
       if d == 4 || d == 7 {
           res++
       }
       x /= 10
   }
   return res
}

func minInt64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

// solve finds minimal a' >= a such that f(a')==x and f(a'+add)==y
func solve(add int64, x, y int, a int64) int64 {
   if f(a) == x && f(a+add) == y {
       return a
   }
   pw := make([]int64, 11)
   pw[0] = 1
   for i := 1; i <= 10; i++ {
       pw[i] = pw[i-1] * 10
   }
   const lnf = int64(9e18)
   // g[pos][countX][countY][carry]
   g := make([][][][]int64, 11)
   for i := 0; i <= 10; i++ {
       g[i] = make([][][]int64, x+1)
       for ca := 0; ca <= x; ca++ {
           g[i][ca] = make([][]int64, y+1)
           for cb := 0; cb <= y; cb++ {
               row := make([]int64, 2)
               row[0], row[1] = lnf, lnf
               g[i][ca][cb] = row
           }
       }
   }
   // initialize
   for d0 := int64(0); d0 < 10; d0++ {
       b := (d0 + add) % 10
       nd := int((d0 + add) / 10)
       if f(d0) <= x && f(b) <= y {
           ca := f(d0)
           cb := f(b)
           if g[0][ca][cb][nd] > d0 {
               g[0][ca][cb][nd] = d0
           }
       }
   }
   // DP over digit positions
   for i := 0; i < 10; i++ {
       for ca := 0; ca <= x; ca++ {
           for cb := 0; cb <= y; cb++ {
               for carry := 0; carry < 2; carry++ {
                   if g[i][ca][cb][carry] < lnf {
                       for na := int64(0); na < 10; na++ {
                           sum := na + int64(carry)
                           nb := sum % 10
                           nc := int(sum / 10)
                           nca := ca + f(na)
                           ncb := cb + f(nb)
                           if nca <= x && ncb <= y {
                               cand := g[i][ca][cb][carry] + na*pw[i+1]
                               if g[i+1][nca][ncb][nc] > cand {
                                   g[i+1][nca][ncb][nc] = cand
                               }
                           }
                       }
                   }
               }
           }
       }
   }
   res := int64(9e18)
   for i := 0; i <= 10; i++ {
       pre := a / pw[i]
       maxAdd := int64(9) - pre%10
       for da := int64(1); da <= maxAdd; da++ {
           if i == 0 {
               if f(pre+da) == x && f(pre+da+add) == y {
                   res = minInt64(res, pre+da)
               }
           } else {
               for carry := 0; carry < 2; carry++ {
                   dx := x - f(pre+da)
                   dy := y - f(pre+da+int64(carry))
                   if dx >= 0 && dy >= 0 && g[i-1][dx][dy][carry] < lnf {
                       cand := (pre+da)*pw[i] + g[i-1][dx][dy][carry]
                       res = minInt64(res, cand)
                   }
               }
           }
       }
   }
   return res
}

// solve10 handles the last up to 10-length segment
func solve10(a, b int64) int64 {
   ans := int64(1e10) + a
   diff := b - a
   for d := 0; d <= 9; d++ {
       fv := make([]int, 20)
       for i := range fv {
           fv[i] = -1000
       }
       for i := int64(0); i <= diff; i++ {
           idx := d + int(i)
           fv[idx] = f(a+i) - f(int64(idx%10))
       }
       pv := make(map[int]struct{})
       sv := make(map[int]struct{})
       for i := 0; i < 10; i++ {
           if fv[i] != -1000 {
               pv[fv[i]] = struct{}{}
           }
           if fv[i+10] != -1000 {
               sv[fv[i+10]] = struct{}{}
           }
       }
       if len(pv) > 1 || len(sv) > 1 {
           continue
       }
       lim := a / 10
       if int64(d) <= a%10 {
           lim++
       }
       if len(sv) == 0 {
           for x := range pv {
               if x >= 0 {
                   cand := solve(0, x, x, lim)*10 + int64(d)
                   ans = minInt64(ans, cand)
               }
           }
       } else {
           var x, y int
           for xv := range pv {
               x = xv
           }
           for yv := range sv {
               y = yv
           }
           if x >= 0 && y >= 0 {
               cand := solve(1, x, y, lim)*10 + int64(d)
               ans = minInt64(ans, cand)
           }
       }
   }
   return ans
}

func main() {
   defer writer.Flush()
   var a, l int64
   fmt.Fscan(reader, &a, &l)
   b := a + l - 1
   base := int64(1)
   ax := int64(0)
   for b >= a+9 {
       ax += (a % 10) * base
       a /= 10
       b /= 10
       base *= 10
   }
   ans := ax + solve10(a, b)*base
   fmt.Fprintln(writer, ans)
}
