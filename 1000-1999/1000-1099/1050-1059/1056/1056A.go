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
   tram := make([]bool, 101)
   for i := 1; i <= 100; i++ {
       tram[i] = true
   }
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       woke := make([]bool, 101)
       for i := 0; i < n; i++ {
           var v int
           fmt.Fscan(reader, &v)
           woke[v] = true
       }
       for i := 1; i <= 100; i++ {
           tram[i] = tram[i] && woke[i]
       }
   }
   first := true
   for i := 1; i <= 100; i++ {
       if tram[i] {
           if !first {
               fmt.Fprint(writer, " ")
           }
           fmt.Fprint(writer, i)
           first = false
       }
   }
   fmt.Fprintln(writer)
}
