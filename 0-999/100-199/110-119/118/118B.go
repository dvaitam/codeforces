package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   total := 2*n + 1
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   for j := 0; j < total; j++ {
       // t is distance from center line
       t := abs(n - j)
       // max digit for this line
       m := n - t
       // print indent spaces
       for i := 0; i < t; i++ {
           writer.WriteByte(' ')
       }
       // print ascending from 0 to m
       for i := 0; i <= m; i++ {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprintf(writer, "%d", i)
       }
       // print descending from m-1 to 0
       for i := m - 1; i >= 0; i-- {
           writer.WriteByte(' ')
           fmt.Fprintf(writer, "%d", i)
       }
       // newline
       writer.WriteByte('\n')
   }
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}
