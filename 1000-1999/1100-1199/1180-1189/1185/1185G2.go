package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, T int
   fmt.Fscan(in, &n, &T)
   var d1, d2, d3 []int
   for i := 0; i < n; i++ {
       var t, g int
       fmt.Fscan(in, &t, &g)
       switch g {
       case 1:
           d1 = append(d1, t)
       case 2:
           d2 = append(d2, t)
       case 3:
           d3 = append(d3, t)
       }
   }
   m1, m2, m3 := len(d1), len(d2), len(d3)
   // DP for group1: F1[k][d]
   dp1 := make([][]int, m1+1)
   for i := range dp1 {
       dp1[i] = make([]int, T+1)
   }
   dp1[0][0] = 1
   for _, t := range d1 {
       for k := m1; k >= 1; k-- {
           rowPrev := dp1[k-1]
           row := dp1[k]
           for d := T; d >= t; d-- {
               row[d] = (row[d] + rowPrev[d-t]) % mod
           }
       }
   }
   // DP combining group1 and group2: dp12[k1][k2][d]
   dp12 := make([][][]int, m1+1)
   for i := 0; i <= m1; i++ {
       dp12[i] = make([][]int, m2+1)
       for j := 0; j <= m2; j++ {
           dp12[i][j] = make([]int, T+1)
       }
       // initialize k2=0
       copy(dp12[i][0], dp1[i])
   }
   for _, t := range d2 {
       for j := m2; j >= 1; j-- {
           for i := 0; i <= m1; i++ {
               prev := dp12[i][j-1]
               cur := dp12[i][j]
               for d := T; d >= t; d-- {
                   cur[d] = (cur[d] + prev[d-t]) % mod
               }
           }
       }
   }
   // DP for group3: dp3[k][d]
   dp3 := make([][]int, m3+1)
   for i := range dp3 {
       dp3[i] = make([]int, T+1)
   }
   dp3[0][0] = 1
   for _, t := range d3 {
       for k := m3; k >= 1; k-- {
           prev := dp3[k-1]
           cur := dp3[k]
           for d := T; d >= t; d-- {
               cur[d] = (cur[d] + prev[d-t]) % mod
           }
       }
   }
   // factorials
   maxm := m1
   if m2 > maxm {
       maxm = m2
   }
   if m3 > maxm {
       maxm = m3
   }
   fact := make([]int, maxm+1)
   fact[0] = 1
   for i := 1; i <= maxm; i++ {
       fact[i] = fact[i-1] * i % mod
   }
   // DP for genre sequences G[a][b][c][last]
   G := make([][][][]int, m1+1)
   for a := 0; a <= m1; a++ {
       G[a] = make([][][]int, m2+1)
       for b := 0; b <= m2; b++ {
           G[a][b] = make([][]int, m3+1)
           for c := 0; c <= m3; c++ {
               G[a][b][c] = make([]int, 4)
           }
       }
   }
   // base cases: single song
   if m1 > 0 {
       G[1][0][0][1] = 1
   }
   if m2 > 0 {
       G[0][1][0][2] = 1
   }
   if m3 > 0 {
       G[0][0][1][3] = 1
   }
   for a := 0; a <= m1; a++ {
       for b := 0; b <= m2; b++ {
           for c := 0; c <= m3; c++ {
               if a+b+c <= 1 {
                   continue
               }
               if a > 0 {
                   v := 0
                   if b > 0 {
                       v = (v + G[a-1][b][c][2]) % mod
                   }
                   if c > 0 {
                       v = (v + G[a-1][b][c][3]) % mod
                   }
                   G[a][b][c][1] = v
               }
               if b > 0 {
                   v := 0
                   if a > 0 {
                       v = (v + G[a][b-1][c][1]) % mod
                   }
                   if c > 0 {
                       v = (v + G[a][b-1][c][3]) % mod
                   }
                   G[a][b][c][2] = v
               }
               if c > 0 {
                   v := 0
                   if a > 0 {
                       v = (v + G[a][b][c-1][1]) % mod
                   }
                   if b > 0 {
                       v = (v + G[a][b][c-1][2]) % mod
                   }
                   G[a][b][c][3] = v
               }
           }
       }
   }
   // compute H[a][b][c]
   H := make([][][]int, m1+1)
   for a := 0; a <= m1; a++ {
       H[a] = make([][]int, m2+1)
       for b := 0; b <= m2; b++ {
           H[a][b] = make([]int, m3+1)
           for c := 0; c <= m3; c++ {
               sum := (G[a][b][c][1] + G[a][b][c][2] + G[a][b][c][3]) % mod
               if sum == 0 {
                   H[a][b][c] = 0
               } else {
                   H[a][b][c] = int(int64(sum) * int64(fact[a]) % mod * int64(fact[b]) % mod * int64(fact[c]) % mod)
               }
           }
       }
   }
   // combine
   var ans int64
   for a := 0; a <= m1; a++ {
       for b := 0; b <= m2; b++ {
           for c := 0; c <= m3; c++ {
               h := H[a][b][c]
               if h == 0 {
                   continue
               }
               for d := 0; d <= T; d++ {
                   cnt12 := dp12[a][b][d]
                   if cnt12 == 0 {
                       continue
                   }
                   d3 := T - d
                   if d3 < 0 || d3 > T {
                       continue
                   }
                   cnt3 := dp3[c][d3]
                   if cnt3 == 0 {
                       continue
                   }
                   ans = (ans + int64(cnt12) * int64(cnt3) % mod * int64(h)) % mod
               }
           }
       }
   }
   fmt.Fprintln(out, ans)
}
