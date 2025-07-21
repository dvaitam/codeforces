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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   prefix := make([]int, n)
   for i := 0; i < n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       if i == 0 {
           prefix[i] = a
       } else {
           prefix[i] = prefix[i-1] + a
       }
   }

   var m int
   fmt.Fscan(reader, &m)
   for i := 0; i < m; i++ {
       var q int
       fmt.Fscan(reader, &q)
       idx := sort.Search(n, func(j int) bool {
           return prefix[j] >= q
       })
       // idx is zero-based, pile numbers are 1-based
       fmt.Fprintln(writer, idx+1)
   }
}
