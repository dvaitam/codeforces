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
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       if n == 2 && m == 2 {
           writer.WriteString("WB\nBB\n")
           continue
       }
       field := make([][]bool, n)
       for i := 0; i < n; i++ {
           field[i] = make([]bool, m)
           for j := 0; j < m; j++ {
               field[i][j] = (j%2 == i%2)
           }
       }
       if n%2 == 0 {
           field[n-1][m-1] = true
           if m >= 2 {
               field[n-1][m-2] = true
           }
       }
       if m%2 == 0 {
           field[n-1][m-1] = true
           if n >= 2 {
               field[n-2][m-1] = true
           }
       }
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               if field[i][j] {
                   writer.WriteByte('B')
               } else {
                   writer.WriteByte('W')
               }
           }
           writer.WriteByte('\n')
       }
   }
}
