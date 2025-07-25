package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var m int64
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   v := make([]int64, n)
   var maxv int64
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &v[i])
       if v[i] > maxv {
           maxv = v[i]
       }
   }
   // Effective max for comparisons
   var M0 int64 = maxv
   if M0 > m {
       M0 = m
   }
   // parity array for diagonal a==b cases
   p := make([]int, M0+1)
   for i := 0; i < n; i++ {
       vi := v[i]
       for d := int64(1); d <= M0; d++ {
           if (vi/d)&1 == 1 {
               p[d] ^= 1
           }
       }
   }
   // counts
   small := m - M0
   // w_s: second player wins (P)
   var wL, wR, wF, wS int64
   // region where no moves for both: (a > M0, b > M0)
   wS = small * small
   // region only Alice moves: a<=M0, b>M0
   wL = M0 * small
   // region only Bob moves: a>M0, b<=M0
   wR = small * M0
   // region both can move: a<=M0, b<=M0
   // for a<b: Alice always wins; a>b: Bob always wins; a==b: diagonal
   // a<b count: M0*(M0-1)/2
   wL += M0 * (M0 - 1) / 2
   wR += M0 * (M0 - 1) / 2
   // diagonal a==b for d in 1..M0
   for d := int64(1); d <= M0; d++ {
       if p[d] == 1 {
           wF++
       } else {
           wS++
       }
   }
   // output: Alice always wins, Bob always wins, first, second
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprintf(out, "%d %d %d %d", wL, wR, wF, wS)
}
