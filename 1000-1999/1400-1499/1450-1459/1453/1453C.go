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

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       grid := make([]string, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &grid[i])
       }
       // Initialize per-digit info
       const D = 10
       rowMinG := make([]int, D)
       rowMaxG := make([]int, D)
       colMinG := make([]int, D)
       colMaxG := make([]int, D)
       rowLeft := make([][]int, D)
       rowRight := make([][]int, D)
       colUp := make([][]int, D)
       colDown := make([][]int, D)
       for d := 0; d < D; d++ {
           rowMinG[d] = n
           colMinG[d] = n
           rowMaxG[d] = -1
           colMaxG[d] = -1
           rowLeft[d] = make([]int, n)
           rowRight[d] = make([]int, n)
           for i := 0; i < n; i++ {
               rowLeft[d][i] = n
               rowRight[d][i] = -1
           }
           colUp[d] = make([]int, n)
           colDown[d] = make([]int, n)
           for j := 0; j < n; j++ {
               colUp[d][j] = n
               colDown[d][j] = -1
           }
       }
       // Scan grid
       for i := 0; i < n; i++ {
           row := grid[i]
           for j := 0; j < n; j++ {
               d := int(row[j] - '0')
               // global
               if i < rowMinG[d] {
                   rowMinG[d] = i
               }
               if i > rowMaxG[d] {
                   rowMaxG[d] = i
               }
               if j < colMinG[d] {
                   colMinG[d] = j
               }
               if j > colMaxG[d] {
                   colMaxG[d] = j
               }
               // per row
               if j < rowLeft[d][i] {
                   rowLeft[d][i] = j
               }
               if j > rowRight[d][i] {
                   rowRight[d][i] = j
               }
               // per col
               if i < colUp[d][j] {
                   colUp[d][j] = i
               }
               if i > colDown[d][j] {
                   colDown[d][j] = i
               }
           }
       }
       // Compute answers
       for d := 0; d < D; d++ {
           ans := 0
           if rowMaxG[d] >= 0 {
               // per row spans
               for i := 0; i < n; i++ {
                   if rowRight[d][i] >= 0 {
                       base := rowRight[d][i] - rowLeft[d][i]
                       height := max(i, n-1-i)
                       ans = max(ans, base*height)
                   }
               }
               // per col spans
               for j := 0; j < n; j++ {
                   if colDown[d][j] >= 0 {
                       base := colDown[d][j] - colUp[d][j]
                       width := max(j, n-1-j)
                       ans = max(ans, base*width)
                   }
               }
               // per cell
               for i := 0; i < n; i++ {
                   row := grid[i]
                   for j := 0; j < n; j++ {
                       if int(row[j]-'0') != d {
                           continue
                       }
                       // horizontal base using two in this row and change at edge
                       heightSpan := max(rowMaxG[d]-i, i-rowMinG[d])
                       distCol := max(j, n-1-j)
                       ans = max(ans, heightSpan*distCol)
                       // vertical base using two in this col and change at edge
                       widthSpan := max(colMaxG[d]-j, j-colMinG[d])
                       distRow := max(i, n-1-i)
                       ans = max(ans, widthSpan*distRow)
                   }
               }
           }
           // print answer
           if d > 0 {
               fmt.Fprint(writer, " ")
           }
           fmt.Fprint(writer, ans)
       }
       fmt.Fprintln(writer)
   }
}
