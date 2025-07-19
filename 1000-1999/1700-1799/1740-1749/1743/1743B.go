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
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for T > 0 {
       T--
       var n int
       fmt.Fscan(reader, &n)
       // print permutation: start with 2
       fmt.Fprint(writer, 2, " ")
       // then 3 to n-1
       for i := 3; i < n; i++ {
           fmt.Fprint(writer, i, " ")
       }
       // then n and 1
       fmt.Fprint(writer, n, " ", 1, "\n")
   }
}
