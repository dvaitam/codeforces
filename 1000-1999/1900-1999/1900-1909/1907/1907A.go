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
   fmt.Fscan(reader, &t)
   columns := "abcdefgh"
   for i := 0; i < t; i++ {
       var s string
       var row int
       fmt.Fscan(reader, &s, &row)
       c := s[0]
       // same row
       for j := 0; j < len(columns); j++ {
           if columns[j] == c {
               continue
           }
           fmt.Fprintf(writer, "%c%d\n", columns[j], row)
       }
       // same column
       for r := 1; r <= 8; r++ {
           if r == row {
               continue
           }
           fmt.Fprintf(writer, "%c%d\n", c, r)
       }
   }
}
