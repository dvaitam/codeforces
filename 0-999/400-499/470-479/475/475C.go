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
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       line := make([]byte, m)
       fmt.Fscan(reader, &line)
       grid[i] = line
   }
   // row segments
   rowL := make([]int, n)
   rowR := make([]int, n)
   for i := 0; i < n; i++ {
       rowL[i] = m
       rowR[i] = -1
       for j := 0; j < m; j++ {
           if grid[i][j] == 'X' {
               if j < rowL[i] {
                   rowL[i] = j
               }
               if j > rowR[i] {
                   rowR[i] = j
               }
           }
       }
   }
   // find first/last row with X
   firstRow, lastRow := -1, -1
   for i := 0; i < n; i++ {
       if rowR[i] >= 0 {
           if firstRow < 0 {
               firstRow = i
           }
           lastRow = i
       }
   }
   // check no empty row inside
   for i := firstRow; i <= lastRow; i++ {
       if rowR[i] < 0 {
           fmt.Println(-1)
           return
       }
       // contiguous check
       if rowR[i]-rowL[i]+1 != count(grid[i], 'X') {
           fmt.Println(-1)
           return
       }
   }
   // minimal width y
   y := m
   for i := firstRow; i <= lastRow; i++ {
       w := rowR[i] - rowL[i] + 1
       if w < y {
           y = w
       }
   }
   // column segments
   colL := make([]int, m)
   colR := make([]int, m)
   for j := 0; j < m; j++ {
       colL[j] = n
       colR[j] = -1
       for i := 0; i < n; i++ {
           if grid[i][j] == 'X' {
               if i < colL[j] {
                   colL[j] = i
               }
               if i > colR[j] {
                   colR[j] = i
               }
           }
       }
   }
   firstCol, lastCol := -1, -1
   for j := 0; j < m; j++ {
       if colR[j] >= 0 {
           if firstCol < 0 {
               firstCol = j
           }
           lastCol = j
       }
   }
   for j := firstCol; j <= lastCol; j++ {
       if colR[j] < 0 {
           fmt.Println(-1)
           return
       }
       // contiguous in column
       if colR[j]-colL[j]+1 != countCol(grid, j, 'X') {
           fmt.Println(-1)
           return
       }
   }
   // minimal height x
   x := n
   for j := firstCol; j <= lastCol; j++ {
       h := colR[j] - colL[j] + 1
       if h < x {
           x = h
       }
   }
   // prefix sum of grid
   sum := make([][]int, n+1)
   for i := range sum {
       sum[i] = make([]int, m+1)
   }
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           sum[i+1][j+1] = sum[i+1][j] + sum[i][j+1] - sum[i][j]
           if grid[i][j] == 'X' {
               sum[i+1][j+1]++
           }
       }
   }
   // erosion: find all brush top-left positions
   A := make([][]bool, n)
   for i := range A {
       A[i] = make([]bool, m)
   }
   for i := 0; i+x <= n; i++ {
       for j := 0; j+y <= m; j++ {
           total := sum[i+x][j+y] - sum[i][j+y] - sum[i+x][j] + sum[i][j]
           if total == x*y {
               A[i][j] = true
           }
       }
   }
   // dilation via diff
   diff := make([][]int, n+1)
   for i := range diff {
       diff[i] = make([]int, m+1)
   }
   for i := 0; i+x <= n; i++ {
       for j := 0; j+y <= m; j++ {
           if A[i][j] {
               diff[i][j]++
               diff[i+x][j]--
               diff[i][j+y]--
               diff[i+x][j+y]++
           }
       }
   }
   // prefix diff to cover
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if i > 0 {
               diff[i][j] += diff[i-1][j]
           }
           if j > 0 {
               diff[i][j] += diff[i][j-1]
           }
           if i > 0 && j > 0 {
               diff[i][j] -= diff[i-1][j-1]
           }
           // check coverage
           if (grid[i][j] == 'X') != (diff[i][j] > 0) {
               fmt.Println(-1)
               return
           }
       }
   }
   fmt.Println(x * y)
}

func count(arr []byte, c byte) int {
   cnt := 0
   for _, v := range arr {
       if v == c {
           cnt++
       }
   }
   return cnt
}

func countCol(grid [][]byte, col int, c byte) int {
   cnt := 0
   for i := range grid {
       if grid[i][col] == c {
           cnt++
       }
   }
   return cnt
}
