package main

import (
   "bufio"
   "fmt"
   "os"
   "math/bits"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(in, &a[i][j])
       }
   }
   W := (m + 63) >> 6
   rows := make([][]uint64, n)
   for i := range rows {
       rows[i] = make([]uint64, W)
   }

   check := func(x int) bool {
       good := make([]int, 0, n)
       for i := 0; i < n; i++ {
           cnt := 0
           for j := 0; j < m; j++ {
               if a[i][j] >= x {
                   w := j >> 6
                   b := uint(j & 63)
                   rows[i][w] |= 1 << b
                   cnt++
               }
           }
           if cnt >= 2 {
               good = append(good, i)
           }
       }
       if len(good) < 2 {
           for i := 0; i < n; i++ {
               for k := range rows[i] {
                   rows[i][k] = 0
               }
           }
           return false
       }
       for ii := 0; ii < len(good); ii++ {
           i := good[ii]
           for jj := ii + 1; jj < len(good); jj++ {
               k := good[jj]
               found := 0
               for w := 0; w < W; w++ {
                   common := rows[i][w] & rows[k][w]
                   if common == 0 {
                       continue
                   }
                   found += bits.OnesCount64(common)
                   if found >= 2 {
                       for t := 0; t < n; t++ {
                           for u := range rows[t] {
                               rows[t][u] = 0
                           }
                       }
                       return true
                   }
               }
           }
       }
       for i := 0; i < n; i++ {
           for k := range rows[i] {
               rows[i][k] = 0
           }
       }
       return false
   }

   lo, hi := 0, int(1e9)
   for lo < hi {
       mid := (lo + hi + 1) >> 1
       if check(mid) {
           lo = mid
       } else {
           hi = mid - 1
       }
   }
   fmt.Fprintln(out, lo)
}
