package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   tab := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &tab[i])
   }
   sorted := make([]int, n)
   copy(sorted, tab)
   sort.Ints(sorted)
   for i, v := range tab {
       for j, sv := range sorted {
           if sv == v {
               fmt.Fprint(writer, sorted[(j+1)%n])
               if i < n-1 {
                   fmt.Fprint(writer, " ")
               }
               break
           }
       }
   }
   fmt.Fprintln(writer)
}
