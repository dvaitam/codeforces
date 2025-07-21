package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func ceilSqrt(x int64) int {
   if x <= 0 {
       return 0
   }
   y := int(math.Sqrt(float64(x)))
   if int64(y)*int64(y) < x {
       y++
   }
   return y
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   var m int64
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   // Count irrational roots
   var irr int64
   for b := 1; b <= n; b++ {
       b2 := int64(b) * int64(b)
       // X = min(m, b^2-1)
       var X int64
       if m < b2-1 {
           X = m
       } else {
           X = b2 - 1
       }
       if X <= 0 {
           continue
       }
       // count D = b2 - c in [1..X] that are perfect squares
       t := b2 - X
       sMin := ceilSqrt(t)
       // number of squares = b - sMin
       sq := b - sMin
       if sq < 0 {
           sq = 0
       }
       irrC := X - int64(sq)
       irr += 2 * irrC
   }
   // Count distinct integer roots via s = a+-d = u or v
   maxS := 2 * n
   mark := make([]byte, maxS+1)
   for u := 1; u <= n; u++ {
       // v >= u, uv <= m, u+v <= 2n, parity same
       vMax1 := int(m / int64(u))
       vMax2 := 2*n - u
       vMax := vMax1
       if vMax2 < vMax {
           vMax = vMax2
       }
       if vMax < u {
           continue
       }
       // mark u
       mark[u] = 1
       // mark v with same parity
       // start from u, step by 2
       start := u
       // ensure start%2 == u%2 always true
       for v := start; v <= vMax; v += 2 {
           mark[v] = 1
       }
   }
   var intCnt int64
   for s := 1; s <= maxS; s++ {
       if mark[s] != 0 {
           intCnt++
       }
   }
   // Total distinct roots
   ans := irr + intCnt
   fmt.Fprintln(out, ans)
}
