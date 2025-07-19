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
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   var c int64
   for i := 0; i+1 < n; i++ {
       diff := a[i+1] - a[i] - 1
       if diff > 0 {
           c += diff
       }
   }
   fmt.Fprintln(writer, c)
}
