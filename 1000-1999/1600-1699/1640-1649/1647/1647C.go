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

   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
       T--
       var n, m int
       fmt.Fscan(reader, &n, &m)
       mat := make([][]byte, n)
       for i := 0; i < n; i++ {
           var line string
           fmt.Fscan(reader, &line)
           mat[i] = []byte(line)
       }
       // count ones
       s := 0
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               if mat[i][j] == '1' {
                   s++
               }
           }
       }
       // impossible if top-left is 1
       if mat[0][0] == '1' {
           fmt.Fprintln(writer, -1)
           continue
       }
       fmt.Fprintln(writer, s)
       // process rows bottom-up except first
       for i := n - 1; i >= 1; i-- {
           for j := 0; j < m; j++ {
               if mat[i][j] == '1' {
                   // connect with cell above
                   fmt.Fprintln(writer, i, j+1, i+1, j+1)
               }
           }
       }
       // process first row leftwards
       for j := m - 1; j >= 1; j-- {
           if mat[0][j] == '1' {
               // connect with cell to the left
               fmt.Fprintln(writer, 1, j, 1, j+1)
           }
       }
   }
}
