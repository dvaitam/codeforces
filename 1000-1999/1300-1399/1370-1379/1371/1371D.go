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
       var n, k int
       fmt.Fscan(reader, &n, &k)
       if k % n != 0 {
           fmt.Fprintln(writer, 2)
       } else {
           fmt.Fprintln(writer, 0)
       }
       // initialize n x n matrix of '0'
       w := make([][]byte, n)
       for i := 0; i < n; i++ {
           row := make([]byte, n)
           for j := 0; j < n; j++ {
               row[j] = '0'
           }
           w[i] = row
       }
       // fill k ones
       x, y := 0, 0
       for i := 0; i < k; i++ {
           w[x][y] = '1'
           x++
           if x == n {
               x = 0
           }
           y++
           if y == n {
               y = 0
               x++
               if x == n {
                   x = 0
               }
           }
       }
       // output matrix
       for i := 0; i < n; i++ {
           writer.Write(w[i])
           writer.WriteByte('\n')
       }
   }
}
