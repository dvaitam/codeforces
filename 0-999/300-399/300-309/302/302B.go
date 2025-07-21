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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   prefix := make([]int64, n)
   var ci, ti int64
   var sum int64
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &ci, &ti)
       sum += ci * ti
       prefix[i] = sum
   }
   // process queries
   for i := 0; i < m; i++ {
       var v int64
       fmt.Fscan(reader, &v)
       idx := sort.Search(len(prefix), func(j int) bool {
           return prefix[j] >= v
       })
       // output song number (1-based)
       fmt.Fprint(writer, idx+1)
       if i+1 < m {
           fmt.Fprint(writer, " ")
       }
   }
   fmt.Fprintln(writer)
}
