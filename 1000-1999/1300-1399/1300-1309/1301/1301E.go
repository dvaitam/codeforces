package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m, q int
   fmt.Fscan(in, &n, &m, &q)
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(in, &s)
       grid[i] = []byte(s)
   }
   // prefix sums for each color: 0:R,1:G,2:Y,3:B
   sums := make([][][]int, 4)
   for c := 0; c < 4; c++ {
       sums[c] = make([][]int, n+1)
       for i := range sums[c] {
           sums[c][i] = make([]int, m+1)
       }
   }
   colorIdx := map[byte]int{'R': 0, 'G': 1, 'Y': 2, 'B': 3}
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           ci := colorIdx[grid[i-1][j-1]]
           for c := 0; c < 4; c++ {
               sums[c][i][j] = sums[c][i-1][j] + sums[c][i][j-1] - sums[c][i-1][j-1]
           }
           sums[ci][i][j]++
       }
   }
   // helper to get sum of color c in rectangle [r1..r2][c1..c2]
   rectSum := func(c, r1, c1, r2, c2 int) int {
       return sums[c][r2][c2] - sums[c][r1-1][c2] - sums[c][r2][c1-1] + sums[c][r1-1][c1-1]
   }
   // bestK[i][j]: max k at 0-based i,j
   bestK := make([][]int, n)
   for i := 0; i < n; i++ {
       bestK[i] = make([]int, m)
       for j := 0; j < m; j++ {
           maxk := min(n-i, m-j) / 2
           lo, hi := 0, maxk
           for lo < hi {
               mid := (lo + hi + 1) >> 1
               // check mid
               k := mid
               r1, c1 := i+1, j+1
               // quadrants:
               if rectSum(0, r1, c1, r1+k-1, c1+k-1) == k*k &&
                   rectSum(1, r1, c1+k, r1+k-1, c1+2*k-1) == k*k &&
                   rectSum(2, r1+k, c1, r1+2*k-1, c1+k-1) == k*k &&
                   rectSum(3, r1+k, c1+k, r1+2*k-1, c1+2*k-1) == k*k {
                   lo = mid
               } else {
                   hi = mid - 1
               }
           }
           bestK[i][j] = lo
       }
   }
   // build logs
   maxNM := max(n, m)
   logs := make([]int, maxNM+1)
   for i := 2; i <= maxNM; i++ {
       logs[i] = logs[i>>1] + 1
   }
   LN := logs[n]
   LM := logs[m]
   // build 2D sparse table st[p][q][i][j]
   st := make([][][][]int, LN+1)
   for p := 0; p <= LN; p++ {
       st[p] = make([][][]int, LM+1)
       for q1 := 0; q1 <= LM; q1++ {
           rows := n - (1 << p) + 1
           cols := m - (1 << q1) + 1
           if rows <= 0 || cols <= 0 {
               continue
           }
           st[p][q1] = make([][]int, rows)
           for i := 0; i < rows; i++ {
               st[p][q1][i] = make([]int, cols)
           }
       }
   }
   // p=0,q=0
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           st[0][0][i][j] = bestK[i][j]
       }
   }
   // p=0, q>0
   for q1 := 1; q1 <= LM; q1++ {
       for i := 0; i < n; i++ {
           for j := 0; j + (1<<q1) <= m; j++ {
               st[0][q1][i][j] = max(st[0][q1-1][i][j], st[0][q1-1][i][j+(1<<(q1-1))])
           }
       }
   }
   // p>0
   for p := 1; p <= LN; p++ {
       for q1 := 0; q1 <= LM; q1++ {
           if st[p][q1] == nil {
               continue
           }
           for i := 0; i + (1<<p) <= n; i++ {
               for j := 0; j + (1<<q1) <= m; j++ {
                   st[p][q1][i][j] = max(st[p-1][q1][i][j], st[p-1][q1][i+(1<<(p-1))][j])
               }
           }
       }
   }
   // function to query maxK in rect [r1..r2][c1..c2], 0-based inclusive
   queryMax := func(r1, c1, r2, c2 int) int {
       if r1 > r2 || c1 > c2 {
           return 0
       }
       h := r2 - r1 + 1
       w := c2 - c1 + 1
       p := logs[h]
       q1 := logs[w]
       x2 := r2 - (1 << p) + 1
       y2 := c2 - (1 << q1) + 1
       m1 := max(st[p][q1][r1][c1], st[p][q1][x2][c1])
       m2 := max(st[p][q1][r1][y2], st[p][q1][x2][y2])
       return max(m1, m2)
   }
   // process queries
   for qi := 0; qi < q; qi++ {
       var r1, c1, r2, c2 int
       fmt.Fscan(in, &r1, &c1, &r2, &c2)
       // convert to 0-based
       r1--
       c1--
       r2--
       c2--
       maxPossible := min(r2-r1+1, c2-c1+1) / 2
       lo, hi := 0, maxPossible
       for lo < hi {
           mid := (lo + hi + 1) >> 1
           // check if any bestK >= mid in starts [r1..r2-2*mid+1][c1..c2-2*mid+1]
           if queryMax(r1, c1, r2-2*mid+1, c2-2*mid+1) >= mid {
               lo = mid
           } else {
               hi = mid - 1
           }
       }
       L := 2 * lo
       fmt.Fprintln(out, L*L)
   }
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}
