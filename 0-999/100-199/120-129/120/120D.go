package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([][]int, n)
   for i := 0; i < n; i++ {
       grid[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &grid[i][j])
       }
   }
   var A, B, C int
   fmt.Fscan(reader, &A, &B, &C)
   target := []int{A, B, C}
   sort.Ints(target)

   // row sums and prefix
   rowSum := make([]int, n)
   for i := 0; i < n; i++ {
       sum := 0
       for j := 0; j < m; j++ {
           sum += grid[i][j]
       }
       rowSum[i] = sum
   }
   prefRow := make([]int, n+1)
   for i := 0; i < n; i++ {
       prefRow[i+1] = prefRow[i] + rowSum[i]
   }

   // column sums and prefix
   colSum := make([]int, m)
   for j := 0; j < m; j++ {
       sum := 0
       for i := 0; i < n; i++ {
           sum += grid[i][j]
       }
       colSum[j] = sum
   }
   prefCol := make([]int, m+1)
   for j := 0; j < m; j++ {
       prefCol[j+1] = prefCol[j] + colSum[j]
   }

   var ways int64
   // horizontal cuts
   if n >= 3 {
       for c1 := 1; c1 <= n-2; c1++ {
           for c2 := c1 + 1; c2 <= n-1; c2++ {
               s1 := prefRow[c1] - prefRow[0]
               s2 := prefRow[c2] - prefRow[c1]
               s3 := prefRow[n] - prefRow[c2]
               parts := []int{s1, s2, s3}
               sort.Ints(parts)
               if parts[0] == target[0] && parts[1] == target[1] && parts[2] == target[2] {
                   ways++
               }
           }
       }
   }
   // vertical cuts
   if m >= 3 {
       for c1 := 1; c1 <= m-2; c1++ {
           for c2 := c1 + 1; c2 <= m-1; c2++ {
               s1 := prefCol[c1] - prefCol[0]
               s2 := prefCol[c2] - prefCol[c1]
               s3 := prefCol[m] - prefCol[c2]
               parts := []int{s1, s2, s3}
               sort.Ints(parts)
               if parts[0] == target[0] && parts[1] == target[1] && parts[2] == target[2] {
                   ways++
               }
           }
       }
   }
   fmt.Println(ways)
}
