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
   // read empty rows
   var t int
   fmt.Fscan(reader, &t)
   emptyRow := make([]bool, n+2)
   for i := 0; i < t; i++ {
       var r int
       fmt.Fscan(reader, &r)
       if r >= 1 && r <= n {
           emptyRow[r] = true
       }
   }
   // read empty columns
   var s int
   fmt.Fscan(reader, &s)
   emptyCol := make([]bool, m+2)
   for i := 0; i < s; i++ {
       var c int
       fmt.Fscan(reader, &c)
       if c >= 1 && c <= m {
           emptyCol[c] = true
       }
   }
   // count segments of non-empty consecutive rows
   rowSeg := 0
   for i := 1; i <= n; i++ {
       if !emptyRow[i] && (i == 1 || emptyRow[i-1]) {
           rowSeg++
       }
   }
   // count segments of non-empty consecutive columns
   colSeg := 0
   for j := 1; j <= m; j++ {
       if !emptyCol[j] && (j == 1 || emptyCol[j-1]) {
           colSeg++
       }
   }
   // total connected components is product of segments
   result := rowSeg * colSeg
   fmt.Println(result)
}
