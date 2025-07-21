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
       var line string
       fmt.Fscan(reader, &line)
       grid[i] = []byte(line)
   }
   // Precompute L, R, U, D
   L := make([][]int, n)
   R := make([][]int, n)
   U := make([][]int, n)
   D := make([][]int, n)
   for i := 0; i < n; i++ {
       L[i] = make([]int, m)
       R[i] = make([]int, m)
       U[i] = make([]int, m)
       D[i] = make([]int, m)
   }
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == '.' {
               if j > 0 {
                   L[i][j] = L[i][j-1] + 1
               } else {
                   L[i][j] = 1
               }
               if i > 0 {
                   U[i][j] = U[i-1][j] + 1
               } else {
                   U[i][j] = 1
               }
           }
       }
       for j := m - 1; j >= 0; j-- {
           if grid[i][j] == '.' {
               if j < m-1 {
                   R[i][j] = R[i][j+1] + 1
               } else {
                   R[i][j] = 1
               }
           }
       }
   }
   for j := 0; j < m; j++ {
       for i := n - 1; i >= 0; i-- {
           if grid[i][j] == '.' {
               if i < n-1 {
                   D[i][j] = D[i+1][j] + 1
               } else {
                   D[i][j] = 1
               }
           }
       }
   }
   var total int64 = 0
   // Straight segments
   // Horizontal: rows 2..n-1 (1..n-2) fully empty
   for i := 1; i < n-1; i++ {
       if L[i][m-1] == m {
           total++
       }
   }
   // Vertical: cols 2..m-1 (1..m-2) fully empty
   for j := 1; j < m-1; j++ {
       if U[n-1][j] == n {
           total++
       }
   }
   // One-turn (L-shape)
   for i := 1; i < n-1; i++ {
       for j := 1; j < m-1; j++ {
           if grid[i][j] != '.' {
               continue
           }
           // left-up
           if L[i][j] >= j+1 && U[i][j] >= i+1 {
               total++
           }
           // left-down
           if L[i][j] >= j+1 && D[i][j] >= n-i {
               total++
           }
           // right-up
           if R[i][j] >= m-j && U[i][j] >= i+1 {
               total++
           }
           // right-down
           if R[i][j] >= m-j && D[i][j] >= n-i {
               total++
           }
       }
   }
   // Two-turn shapes: left-to-right via vertical
   for j := 1; j < m-1; j++ {
       var sumA int64 = 0
       for i := 1; i < n-1; i++ {
           if grid[i][j] == '.' {
               // b: right arm
               if R[i][j] >= m-j {
                   total += sumA
               }
               // a: left arm
               if L[i][j] >= j+1 {
                   sumA++
               }
           } else {
               sumA = 0
           }
       }
   }
   // Two-turn shapes: top-to-bottom via horizontal
   for i := 1; i < n-1; i++ {
       var sumA int64 = 0
       for j := 1; j < m-1; j++ {
           if grid[i][j] == '.' {
               // b: down arm
               if D[i][j] >= n-i {
                   total += sumA
               }
               // a: up arm
               if U[i][j] >= i+1 {
                   sumA++
               }
           } else {
               sumA = 0
           }
       }
   }
   // Output result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, total)
}
