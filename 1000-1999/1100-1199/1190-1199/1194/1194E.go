package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   type HSeg struct{ x1, x2, y int }
   type VSeg struct{ x, y1, y2 int }
   hs := make([]HSeg, 0, n)
   vs := make([]VSeg, 0, n)

   for i := 0; i < n; i++ {
       var x1, y1, x2, y2 int
       fmt.Fscan(in, &x1, &y1, &x2, &y2)
       if y1 == y2 {
           // horizontal
           if x1 > x2 {
               x1, x2 = x2, x1
           }
           hs = append(hs, HSeg{x1, x2, y1})
       } else {
           // vertical
           if y1 > y2 {
               y1, y2 = y2, y1
           }
           vs = append(vs, VSeg{x1, y1, y2})
       }
   }
   H := len(hs)
   V := len(vs)
   if H < 2 || V < 2 {
       fmt.Fprintln(out, 0)
       return
   }
   // build bitsets: for each vertical, bit j if it intersects horizontal j
   w := (H + 63) >> 6
   bitsArr := make([][]uint64, V)
   for i := 0; i < V; i++ {
       bitsArr[i] = make([]uint64, w)
       vx := vs[i].x
       y1, y2 := vs[i].y1, vs[i].y2
       for j := 0; j < H; j++ {
           h := hs[j]
           if h.y >= y1 && h.y <= y2 && vx >= h.x1 && vx <= h.x2 {
               bitsArr[i][j>>6] |= 1 << (uint(j) & 63)
           }
       }
   }
   var ans int64
   // for each pair of verticals
   for i := 0; i < V; i++ {
       for j := i + 1; j < V; j++ {
           // count horizontals crossing both
           var cnt int64
           ai, aj := bitsArr[i], bitsArr[j]
           for k := 0; k < w; k++ {
               cnt += int64(bits.OnesCount64(ai[k] & aj[k]))
           }
           if cnt > 1 {
               ans += cnt * (cnt - 1) / 2
           }
       }
   }
   fmt.Fprintln(out, ans)
}
