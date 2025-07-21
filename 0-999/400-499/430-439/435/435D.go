package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([][]byte, n+2)
   for i := 1; i <= n; i++ {
       line := make([]byte, m+2)
       var s string
       fmt.Fscan(reader, &s)
       for j := 1; j <= m; j++ {
           line[j] = s[j-1]
       }
       grid[i] = line
   }
   // black indicator: 1 if black, 0 if white
   black := make([][]int, n+2)
   for i := range black {
       black[i] = make([]int, m+2)
   }
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           if grid[i][j] == '1' {
               black[i][j] = 1
           }
       }
   }
   // prefix sums
   // row and col
   rowPS := make([][]int, n+2)
   colPS := make([][]int, n+2)
   for i := range rowPS {
       rowPS[i] = make([]int, m+2)
       colPS[i] = make([]int, m+2)
   }
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           rowPS[i][j] = rowPS[i][j-1] + black[i][j]
           colPS[i][j] = colPS[i-1][j] + black[i][j]
       }
   }
   // diag1 (slope 1) and diag2 (slope -1)
   diag1 := make([][]int, n+2)
   diag2 := make([][]int, n+2)
   for i := range diag1 {
       diag1[i] = make([]int, m+3)
       diag2[i] = make([]int, m+3)
   }
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           diag1[i][j] = black[i][j] + diag1[i-1][j-1]
       }
       for j := m; j >= 1; j-- {
           diag2[i][j] = black[i][j] + diag2[i-1][j+1]
       }
   }
   // run lengths of white in four directions
   hRun := make([][]int, n+2)    // right
   hRunL := make([][]int, n+2)   // left
   vRun := make([][]int, n+2)    // down
   vRunU := make([][]int, n+2)   // up
   for i := range hRun {
       hRun[i] = make([]int, m+2)
       hRunL[i] = make([]int, m+2)
       vRun[i] = make([]int, m+2)
       vRunU[i] = make([]int, m+2)
   }
   for i := n; i >= 1; i-- {
       for j := m; j >= 1; j-- {
           if grid[i][j] == '0' {
               hRun[i][j] = hRun[i][j+1] + 1
               vRun[i][j] = vRun[i+1][j] + 1
           }
       }
   }
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           if grid[i][j] == '0' {
               hRunL[i][j] = hRunL[i][j-1] + 1
               vRunU[i][j] = vRunU[i-1][j] + 1
           }
       }
   }
   var ans int64
   // Type1: right isosceles with axis legs, four orientations
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           if grid[i][j] != '0' {
               continue
           }
           // right & down
           maxK := hRun[i][j]
           if vRun[i][j] < maxK {
               maxK = vRun[i][j]
           }
           for k := 1; k < maxK; k++ {
               // endpoints
               if grid[i][j+k] != '0' || grid[i+k][j] != '0' {
                   continue
               }
               if rowPS[i][j+k-1]-rowPS[i][j] != 0 {
                   continue
               }
               if colPS[i+k-1][j]-colPS[i][j] != 0 {
                   continue
               }
               // diag slope -1 from B to C
               if diag2[i+k-1][j+1] - diag2[i][j+k] != 0 {
                   continue
               }
               ans++
           }
           // right & up
           maxK = hRun[i][j]
           if vRunU[i][j] < maxK {
               maxK = vRunU[i][j]
           }
           for k := 1; k < maxK; k++ {
               if grid[i][j+k] != '0' || grid[i-k][j] != '0' {
                   continue
               }
               if rowPS[i][j+k-1]-rowPS[i][j] != 0 {
                   continue
               }
               if colPS[i-1][j]-colPS[i-k][j] != 0 {
                   continue
               }
               // diag slope 1 from C to B
               if diag1[i-1][j+k-1] - diag1[i-k][j] != 0 {
                   continue
               }
               ans++
           }
           // left & down
           maxK = hRunL[i][j]
           if vRun[i][j] < maxK {
               maxK = vRun[i][j]
           }
           for k := 1; k < maxK; k++ {
               if grid[i][j-k] != '0' || grid[i+k][j] != '0' {
                   continue
               }
               if rowPS[i][j-1]-rowPS[i][j-k] != 0 {
                   continue
               }
               if colPS[i+k-1][j]-colPS[i][j] != 0 {
                   continue
               }
               // diag slope 1 from B to C
               if diag1[i+k-1][j-1] - diag1[i][j-k] != 0 {
                   continue
               }
               ans++
           }
           // left & up
           maxK = hRunL[i][j]
           if vRunU[i][j] < maxK {
               maxK = vRunU[i][j]
           }
           for k := 1; k < maxK; k++ {
               if grid[i][j-k] != '0' || grid[i-k][j] != '0' {
                   continue
               }
               if rowPS[i][j-1]-rowPS[i][j-k] != 0 {
                   continue
               }
               if colPS[i-1][j]-colPS[i-k][j] != 0 {
                   continue
               }
               // diag slope -1 from B to C
               if diag2[i-1][j-k+1] - diag2[i-k][j] != 0 {
                   continue
               }
               ans++
           }
       }
   }
   // Type2: horizontal base, two diag sides
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           if grid[i][j] != '0' {
               continue
           }
           // max horizontal length
           maxL := hRun[i][j]
           // base length k even
           for k := 2; k < maxL; k += 2 {
               mid := k / 2
               bj := j + k
               if bj > m || grid[i][bj] != '0' {
                   continue
               }
               // base clear
               if rowPS[i][bj-1]-rowPS[i][j] != 0 {
                   continue
               }
               // above
               ci := i - mid
               cj := j + mid
               if ci >= 1 {
                   if grid[ci][cj] == '0' {
                       // A->C slope -1
                       if diag2[i-1][j+1] - diag2[ci][cj] == 0 &&
                           // B->C slope 1
                           diag1[i-1][bj-1] - diag1[ci][cj] == 0 {
                           ans++
                       }
                   }
               }
               // below
               ci = i + mid
               if ci <= n {
                   if grid[ci][cj] == '0' {
                       // A->C slope 1
                       if diag1[ci-1][cj-1] - diag1[i][j] == 0 &&
                           // B->C slope -1
                           diag2[ci-1][cj+1] - diag2[i][bj] == 0 {
                           ans++
                       }
                   }
               }
           }
       }
   }
   // Type3: vertical base, two diag sides
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           if grid[i][j] != '0' {
               continue
           }
           maxL := vRun[i][j]
           for k := 2; k < maxL; k += 2 {
               mid := k / 2
               bi := i + k
               if bi > n || grid[bi][j] != '0' {
                   continue
               }
               // base clear
               if colPS[bi-1][j]-colPS[i][j] != 0 {
                   continue
               }
               // left
               ci := i + mid
               cj := j - mid
               if cj >= 1 && grid[ci][cj] == '0' {
                   // A->C slope -1
                   if diag2[ci-1][cj+1] - diag2[i][j] == 0 &&
                       // B->C slope 1
                       diag1[bi-1][j-1] - diag1[ci][cj] == 0 {
                       ans++
                   }
               }
               // right
               cj = j + mid
               if cj <= m && grid[ci][cj] == '0' {
                   // A->C slope 1
                   if diag1[ci-1][cj-1] - diag1[i][j] == 0 &&
                       // B->C slope -1
                       diag2[bi-1][j+1] - diag2[ci][cj] == 0 {
                       ans++
                   }
               }
           }
       }
   }
   fmt.Println(ans)
}
