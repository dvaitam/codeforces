package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i, v := range a {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
