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

   var n int
   fmt.Fscan(reader, &n)
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, n)
       for j := 0; j < n; j++ {
           fmt.Fscan(reader, &a[i][j])
       }
   }

   good := make([]int, 0)
   for i := 0; i < n; i++ {
       ok := true
       for j := 0; j < n; j++ {
           if i == j {
               continue
           }
           if a[i][j] == 1 || a[i][j] == 3 {
               ok = false
               break
           }
       }
       if ok {
           good = append(good, i+1)
       }
   }

   fmt.Fprintln(writer, len(good))
   for idx, v := range good {
       if idx > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
