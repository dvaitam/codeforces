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
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)
   // Greedy match: for each original value, find a strictly greater value
   cnt := 0
   j := 0
   for i := 0; i < n; i++ {
       // advance j until a[j] > a[i]
       for j < n && a[j] <= a[i] {
           j++
       }
       if j < n {
           cnt++
           j++
       } else {
           break
       }
   }
   fmt.Fprintln(writer, cnt)
}
