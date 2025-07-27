package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func maxInt64(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   var k int64
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       if err == io.EOF {
           return
       }
       panic(err)
   }
   a := make([]int64, n+1)
   b := make([]int64, m+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   for j := 1; j <= m; j++ {
       fmt.Fscan(in, &b[j])
   }
   // read c
   c := make([][]int64, n+1)
   for i := 1; i <= n; i++ {
       c[i] = make([]int64, m+1)
       for j := 1; j <= m; j++ {
           fmt.Fscan(in, &c[i][j])
       }
   }
   // prefix max and indices
   origG := make([][]int64, n+1)
   idxX := make([][]int, n+1)
   idxY := make([][]int, n+1)
   for i := 0; i <= n; i++ {
       origG[i] = make([]int64, m+1)
       idxX[i] = make([]int, m+1)
       idxY[i] = make([]int, m+1)
   }
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           // start with c[i][j]
           v := c[i][j]
           x, y := i, j
           // from top
           if origG[i-1][j] > v {
               v = origG[i-1][j]
               x = idxX[i-1][j]
               y = idxY[i-1][j]
           }
           // from left
           if origG[i][j-1] > v {
               v = origG[i][j-1]
               x = idxX[i][j-1]
               y = idxY[i][j-1]
           }
           origG[i][j] = v
           idxX[i][j] = x
           idxY[i][j] = y
       }
   }
   // initial unlocked
   var A0, B0 int
   for i := 1; i <= n; i++ {
       if a[i] <= 0 {
           A0 = i
       }
   }
   for j := 1; j <= m; j++ {
       if b[j] <= 0 {
           B0 = j
       }
   }
   type state struct{A, B int}
   var states []state
   // simulate unboosted
   var S int64
   A, B := A0, B0
   runs0 := 0
   for A < n || B < m {
       var nextT int64 = 1<<62
       if A < n && a[A+1] < nextT {
           nextT = a[A+1]
       }
       if B < m && b[B+1] < nextT {
           nextT = b[B+1]
       }
       g := origG[A][B]
       need := nextT - S
       var t int
       if need <= 0 {
           t = 0
       } else {
           t = int((need + g - 1) / g)
       }
       runs0 += t
       S += int64(t) * g
       states = append(states, state{A, B})
       // unlock
       for A < n && a[A+1] <= S {
           A++
       }
       for B < m && b[B+1] <= S {
           B++
       }
   }
   ans := runs0
   // unique candidates
   cand := make(map[int]map[int]bool)
   for _, st := range states {
       x := idxX[st.A][st.B]
       y := idxY[st.A][st.B]
       if cand[x] == nil {
           cand[x] = make(map[int]bool)
       }
       cand[x][y] = true
   }
   // test each candidate
   for x, row := range cand {
       for y := range row {
           // if no effect, skip
           if k == 0 {
               continue
           }
           // simulate
           S2 := int64(0)
           A2, B2 := A0, B0
           runs := 0
           for A2 < n || B2 < m {
               var nextT int64 = 1<<62
               if A2 < n && a[A2+1] < nextT {
                   nextT = a[A2+1]
               }
               if B2 < m && b[B2+1] < nextT {
                   nextT = b[B2+1]
               }
               g := origG[A2][B2]
               if A2 >= x && B2 >= y {
                   boosted := c[x][y] + k
                   if boosted > g {
                       g = boosted
                   }
               }
               need := nextT - S2
               var t int
               if need <= 0 {
                   t = 0
               } else {
                   t = int((need + g - 1) / g)
               }
               runs += t
               if runs >= ans {
                   break
               }
               S2 += int64(t) * g
               for A2 < n && a[A2+1] <= S2 {
                   A2++
               }
               for B2 < m && b[B2+1] <= S2 {
                   B2++
               }
           }
           if runs < ans {
               ans = runs
           }
       }
   }
   fmt.Println(ans)
}
