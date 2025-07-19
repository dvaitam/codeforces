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
       var n int
       fmt.Fscan(reader, &n)
       s := 1
       l := n * n
       // build and print each row
       for i := 0; i < n; i++ {
           row := make([]int, n)
           if i%2 == 0 {
               for j := 0; j < n; j++ {
                   if j%2 == 0 {
                       row[j] = s
                       s++
                   } else {
                       row[j] = l
                       l--
                   }
               }
           } else {
               for j := n - 1; j >= 0; j-- {
                   if j%2 == 0 {
                       row[j] = l
                       l--
                   } else {
                       row[j] = s
                       s++
                   }
               }
           }
           // print row
           for j := 0; j < n; j++ {
               fmt.Fprint(writer, row[j])
               if j+1 < n {
                   writer.WriteByte(' ')
               }
           }
           writer.WriteByte('\n')
       }
       // blank line after each case
       writer.WriteByte('\n')
   }
}
