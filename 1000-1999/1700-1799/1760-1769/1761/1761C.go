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
       // prepare result slices
       res := make([][]int, n)
       for i := 0; i < n; i++ {
           // each set contains itself
           res[i] = []int{i + 1}
       }
       // read matrix rows
       for i := 0; i < n; i++ {
           var s string
           fmt.Fscan(reader, &s)
           for j := 0; j < n && j < len(s); j++ {
               if s[j] == '1' {
                   // add row index (i+1) to column j
                   res[j] = append(res[j], i+1)
               }
           }
       }
       // output
       for i := 0; i < n; i++ {
           row := res[i]
           fmt.Fprint(writer, len(row))
           for _, v := range row {
               fmt.Fprint(writer, " ", v)
           }
           fmt.Fprint(writer, '\n')
       }
   }
}
