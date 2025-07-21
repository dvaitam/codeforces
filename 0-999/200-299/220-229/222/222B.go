package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   // Read initial table
   table := make([][]int, n)
   for i := 0; i < n; i++ {
       row := make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &row[j])
       }
       table[i] = row
   }
   // Initialize row and column mappings
   rowIdx := make([]int, n)
   for i := 0; i < n; i++ {
       rowIdx[i] = i
   }
   colIdx := make([]int, m)
   for j := 0; j < m; j++ {
       colIdx[j] = j
   }

   // Process queries
   for qi := 0; qi < k; qi++ {
       var op byte
       var x, y int
       // Read operation and indices
       fmt.Fscan(reader, &op, &x, &y)
       switch op {
       case 'r':
           // swap rows x and y
           rowIdx[x-1], rowIdx[y-1] = rowIdx[y-1], rowIdx[x-1]
       case 'c':
           // swap columns x and y
           colIdx[x-1], colIdx[y-1] = colIdx[y-1], colIdx[x-1]
       case 'g':
           // get value at (x,y)
           r := rowIdx[x-1]
           c := colIdx[y-1]
           fmt.Fprintln(writer, table[r][c])
       }
   }
}
