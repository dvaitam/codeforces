package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for t > 0 {
       t--
       var n, m int
       fmt.Fscan(reader, &n, &m)
       // Read rows and map element to row index
       rowLists := make([][]int, n)
       elemToRow := make([]int, n*m+1)
       for i := 0; i < n; i++ {
           row := make([]int, m)
           for j := 0; j < m; j++ {
               fmt.Fscan(reader, &row[j])
               elemToRow[row[j]] = i
           }
           rowLists[i] = row
       }
       // Read columns, pick the first as reference for row order
       var firstCol []int
       for j := 0; j < m; j++ {
           col := make([]int, n)
           for i := 0; i < n; i++ {
               fmt.Fscan(reader, &col[i])
           }
           if j == 0 {
               firstCol = col
           }
       }
       // Determine original row order
       rowsOrder := make([]int, n)
       for i := 0; i < n; i++ {
           rowsOrder[i] = elemToRow[firstCol[i]]
       }
       // Output the matrix in restored order
       for _, rid := range rowsOrder {
           row := rowLists[rid]
           for j, v := range row {
               if j > 0 {
                   writer.WriteByte(' ')
               }
               writer.WriteString(strconv.Itoa(v))
           }
           writer.WriteByte('\n')
       }
   }
}
