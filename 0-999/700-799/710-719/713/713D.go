package main

import (
   "bufio"
   "fmt"
   "os"
)

// St stores a 2D sparse table block: arr is flat of size n1*m1, row width m1
type St struct {
   n1, m1 int
   arr    []uint16
}

func maxu(a, b uint16) uint16 {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   a := make([][]uint16, n)
   for i := 0; i < n; i++ {
       a[i] = make([]uint16, m)
       for j := 0; j < m; j++ {
           var x int
           fmt.Fscan(reader, &x)
           if x != 0 {
               a[i][j] = 1
           }
       }
   }
   // dp: max square ending at i,j
   dp := make([][]uint16, n)
   for i := range dp {
       dp[i] = make([]uint16, m)
   }
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if a[i][j] > 0 {
               if i > 0 && j > 0 {
                   v := dp[i-1][j]
                   if dp[i][j-1] < v {
                       v = dp[i][j-1]
                   }
                   if dp[i-1][j-1] < v {
                       v = dp[i-1][j-1]
                   }
                   dp[i][j] = v + 1
               } else {
                   dp[i][j] = 1
               }
           }
       }
   }
   // logs
   maxN := n
   if m > maxN {
       maxN = m
   }
   log2 := make([]int, maxN+1)
   log2[1] = 0
   for i := 2; i <= maxN; i++ {
       log2[i] = log2[i/2] + 1
   }
   K1 := log2[n]
   K2 := log2[m]
   // build sparse table st[k1][k2]
   st := make([][]St, K1+1)
   for k1 := 0; k1 <= K1; k1++ {
       st[k1] = make([]St, K2+1)
   }
   // k1=0,k2=0
   st00 := &st[0][0]
   st00.n1 = n
   st00.m1 = m
   st00.arr = make([]uint16, n*m)
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           st00.arr[i*m+j] = dp[i][j]
       }
   }
   // k1=0, k2>0
   for k2 := 1; k2 <= K2; k2++ {
       prev := &st[0][k2-1]
       curr := &st[0][k2]
       span := 1 << (k2 - 1)
       curr.n1 = n
       curr.m1 = m - (1 << k2) + 1
       curr.arr = make([]uint16, curr.n1*curr.m1)
       for i := 0; i < n; i++ {
           for j := 0; j < curr.m1; j++ {
               v1 := prev.arr[i*prev.m1+j]
               v2 := prev.arr[i*prev.m1+j+span]
               curr.arr[i*curr.m1+j] = maxu(v1, v2)
           }
       }
   }
   // k1>0
   for k1 := 1; k1 <= K1; k1++ {
       span1 := 1 << (k1 - 1)
       for k2 := 0; k2 <= K2; k2++ {
           above := &st[k1-1][k2]
           curr := &st[k1][k2]
           curr.n1 = n - (1 << k1) + 1
           curr.m1 = m - (1 << k2) + 1
           curr.arr = make([]uint16, curr.n1*curr.m1)
           for i := 0; i < curr.n1; i++ {
               for j := 0; j < curr.m1; j++ {
                   v1 := above.arr[i*above.m1+j]
                   v2 := above.arr[(i+span1)*above.m1+j]
                   curr.arr[i*curr.m1+j] = maxu(v1, v2)
               }
           }
       }
   }
   // process queries
   var t int
   fmt.Fscan(reader, &t)
   for qi := 0; qi < t; qi++ {
       var x1, y1, x2, y2 int
       fmt.Fscan(reader, &x1, &y1, &x2, &y2)
       x1--
       y1--
       x2--
       y2--
       lo, hi := 1, min(x2-x1+1, y2-y1+1)
       ans := 0
       for lo <= hi {
           mid := (lo + hi) >> 1
           bx := x1 + mid - 1
           by := y1 + mid - 1
           if bx > x2 || by > y2 {
               hi = mid - 1
               continue
           }
           // query max in [bx..x2][by..y2]
           dx := x2 - bx + 1
           dy := y2 - by + 1
           kx := log2[dx]
           ky := log2[dy]
           s := &st[kx][ky]
           // four corners
           i1, j1 := bx, by
           i2, j2 := x2-(1<<kx)+1, by
           i3, j3 := bx, y2-(1<<ky)+1
           i4, j4 := i2, j3
           m1 := s.arr[i1*s.m1+j1]
           m2 := s.arr[i2*s.m1+j2]
           m3 := s.arr[i3*s.m1+j3]
           m4 := s.arr[i4*s.m1+j4]
           mv := maxu(maxu(m1, m2), maxu(m3, m4))
           if int(mv) >= mid {
               ans = mid
               lo = mid + 1
           } else {
               hi = mid - 1
           }
       }
       fmt.Fprintln(writer, ans)
   }
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}
